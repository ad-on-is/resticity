package internal

import (
	"context"
	"embed"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

type Job struct {
	Id       string `json:"id"`
	job      gocron.Job
	Schedule Schedule `json:"schedule"`
	Running  bool     `json:"running"`
	Force    bool     `json:"force"`
	Ctx      context.Context
	Cancel   context.CancelFunc
}

type Scheduler struct {
	Gocron   gocron.Scheduler
	restic   *Restic
	Jobs     []Job
	jmu      sync.Mutex
	settings *Settings
	OutputCh *chan ChanMsg
	ErrorCh  *chan ChanMsg
	Assets   *embed.FS
}

func NewScheduler(
	settings *Settings,
	restic *Restic,
	outch *chan ChanMsg,
	errch *chan ChanMsg,
) (*Scheduler, error) {

	s := &Scheduler{}
	s.settings = settings
	s.restic = restic
	s.OutputCh = outch
	s.ErrorCh = errch

	if gc, err := gocron.NewScheduler(); err == nil {
		s.Gocron = gc
		s.Gocron.Start()
		return s, nil
	} else {
		return nil, err
	}

}

func (s *Scheduler) RunJobById(id string) {
	for i, j := range s.Jobs {
		if j.Id == id {
			log.Info("Running job manually", "id", id)
			s.Jobs[i].Force = true

			if err := j.job.RunNow(); err != nil {
				log.Error("Error running job manually", "id", id, "err", err)
			}
			break
		}
	}
}

func (s *Scheduler) StopJobById(id string) {
	for _, j := range s.Jobs {
		if j.Id == id {
			(*s.OutputCh) <- ChanMsg{Id: j.Schedule.Id, Msg: "{\"running\": false}", Time: time.Now()}
			j.Cancel()
			log.Warn("Canceling context", "id", id)
			break
		}
	}
}

func (s *Scheduler) DeleteRunningJob(id string) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.Jobs {
		if j.Id == id {

			log.Debug("Stopping running job", "id", id)
			s.Jobs[i].Running = false
			s.Jobs[i].Force = false
			break
		}
	}
}

func (s *Scheduler) FindJobById(id string) *Job {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for _, j := range s.Jobs {
		if j.Id == id {
			return &j
		}
	}
	return nil
}

func (s *Scheduler) SetRunningJob(id string) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.Jobs {
		if j.Id == id {

			s.Jobs[i].Running = true
			log.Debug("Setting forced running job", "id", id)

			break
		}
	}
}

func (s *Scheduler) RecreateCtx(name string) {
	for i, j := range s.Jobs {
		if j.job.Name() == name {
			log.Debug("Recreating context for job", "id", name)
			ctx, cancel := context.WithCancel(context.Background())
			s.Jobs[i].Ctx = ctx
			s.Jobs[i].Cancel = cancel
			break
		}
	}
}

func (s *Scheduler) GetRunningJobs() []Job {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	return funk.Filter(s.Jobs, func(j Job) bool { return j.Running == true }).([]Job)
}

func (s *Scheduler) Notifiy(schedule Schedule, finished bool, hasError bool) {
	what := ""
	from := ""
	to := ""
	switch schedule.Action {
	case "backup":
		what = "Backup"
		break
	case "copy-snapshots":
		what = "Copy snapshots"
		break
	case "prune-repository":
		what = "Prune repository"
	}
	if schedule.FromRepositoryId != "" {
		r := s.settings.Config.GetRepositoryById(schedule.FromRepositoryId)
		from = r.Name
	}

	if schedule.BackupId != "" {
		b := s.settings.Config.GetBackupById(schedule.BackupId)
		from = b.Name
	}
	if schedule.ToRepositoryId != "" {
		r := s.settings.Config.GetRepositoryById(schedule.ToRepositoryId)
		to = r.Name
	}
	action := "started"
	if finished {
		action = "finished"
	}
	title := fmt.Sprintf("%s %s", what, action)
	description := fmt.Sprintf("From %s to %s", from, to)
	if schedule.Action == "prune-repository" {
		description = fmt.Sprintf("On %s", to)
	}
	if hasError {
		title += " with error"
	}
	beeep.Notify(title, description, "")
}

func (s *Scheduler) RescheduleBackups() {

	running := s.GetRunningJobs()
	log.Debug("Terminating running jobs", "jobs", len(running))
	for _, j := range running {
		s.StopJobById(j.Id)
	}

	s.Jobs = []Job{}
	log.Info("Rescheduling backups")

	s.settings.Refresh()
	config := s.settings.Config

	for i := range config.Schedules {
		schedule := config.Schedules[i]
		t := time.Now().AddDate(1000, 0, 0)
		jobDef := gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(t))

		if schedule.Cron != "" {
			jobDef = gocron.CronJob(schedule.Cron, false)
		}

		j, err := s.Gocron.NewJob(
			jobDef,
			gocron.NewTask(func() {

				s.restic.RunSchedule(s.FindJobById(schedule.Id))

			}),
			gocron.WithName(schedule.Id),
			gocron.WithTags(
				"backup:"+schedule.BackupId,
				"repository:"+schedule.ToRepositoryId,
				"fromrepository:"+schedule.ToRepositoryId,
			),
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {

					(*s.OutputCh) <- ChanMsg{Id: jobName, Msg: "{\"running\": true}", Time: time.Now()}

					log.Debug(
						"before job run",
						"id",
						jobName,
					)
					s.SetRunningJob(jobName)
					if config.AppSettings.Notifications.OnScheduleStart {
						s.Notifiy(schedule, false, false)
					}
				}),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {

						(*s.OutputCh) <- ChanMsg{Id: jobName, Msg: "{\"running\": false}", Time: time.Now()}

						log.Debug("after job run", "res", "success", "id", jobName)

						if config.AppSettings.Notifications.OnScheduleSuccess {
							s.Notifiy(schedule, true, false)
						}
						s.DeleteRunningJob(jobName)
						s.RecreateCtx(jobName)
						s.settings.SetLastRun(jobName, "")

					},
				),
				gocron.AfterJobRunsWithError(
					func(jobID uuid.UUID, jobName string, err error) {

						(*s.OutputCh) <- ChanMsg{Id: jobName, Msg: "{\"running\": false}", Time: time.Now()}

						if config.AppSettings.Notifications.OnScheduleError {
							s.Notifiy(schedule, true, true)
						}
						s.DeleteRunningJob(jobName)
						s.RecreateCtx(jobName)
						log.Warn("after job run", "res", "error", "id", jobName, "err", err)
						s.settings.SetLastRun(jobName, err.Error())
					},
				),
			))

		if err != nil {
			log.Error("Error creating Job", "err", err)
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())

		s.Jobs = append(
			s.Jobs,
			Job{
				job:      j,
				Schedule: schedule,
				Id:       schedule.Id,
				Running:  false,
				Force:    false,
				Ctx:      ctx,
				Cancel:   cancel,
			},
		)

	}

	log.Debug("Rerunning terminated jobs", "jobs", len(running))

	for _, r := range running {

		for _, j := range s.Jobs {
			if j.Id == r.Id {
				time.Sleep(1 * time.Second)
				s.RunJobById(j.Id)
			}
		}
	}

}
