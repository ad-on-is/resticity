package internal

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
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
	Chan     chan string
}

type Scheduler struct {
	Gocron   gocron.Scheduler
	restic   *Restic
	Jobs     []Job
	jmu      sync.Mutex
	settings *Settings
	outputCh *chan ChanMsg
}

func NewScheduler(settings *Settings, restic *Restic, ch *chan ChanMsg) (*Scheduler, error) {

	s := &Scheduler{}
	s.settings = settings
	s.restic = restic
	s.outputCh = ch
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
			log.Debug("Running job manually", "id", id)
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
			j.Cancel()
			break
		}
	}
}

func (s *Scheduler) DeleteRunningJob(id string) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.Jobs {
		if j.Id == id {
			go func() {
				if j.Chan != nil {
					j.Chan <- "{\"running\": false}"
				}
			}()
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
			log.Debug("Setting forced running job", "id", id)
			s.Jobs[i].Running = true
			go func() {
				if j.Chan != nil {
					j.Chan <- "{\"running\": true}"
				}
			}()
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

func (s *Scheduler) RescheduleBackups() {

	s.Jobs = []Job{}
	log.Info("Rescheduling backups")

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

					log.Debug("before job run", "id", jobName)
					s.SetRunningJob(jobName)
				}),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {
						log.Debug("after job run", "res", "success", "id", jobName)
						s.DeleteRunningJob(jobName)
						s.RecreateCtx(jobName)
						s.settings.SetLastRun(jobName, "")
					},
				),
				gocron.AfterJobRunsWithError(
					func(jobID uuid.UUID, jobName string, err error) {
						log.Debug("after job run", "res", "error", "id", jobName, "err", err)
						s.DeleteRunningJob(jobName)
						s.RecreateCtx(jobName)
						s.settings.SetLastRun(jobName, err.Error())
					},
				),
			))

		if err != nil {
			log.Error("Error creating Job", "err", err)
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan string)
		go func(c chan string) {
			for msg := range c {
				*s.outputCh <- ChanMsg{Id: schedule.Id, Out: msg, Schedule: schedule}
			}
		}(ch)
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
				Chan:     ch,
			},
		)

	}

}
