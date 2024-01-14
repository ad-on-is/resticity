package main

import (
	"context"
	"fmt"
	"sync"

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
	gocron   gocron.Scheduler
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
		s.gocron = gc
		s.gocron.Start()
		return s, nil
	} else {
		return nil, err
	}

}

func (s *Scheduler) RunJobById(id string) {
	fmt.Println("should run", id)
	for i, j := range s.Jobs {
		if j.Id == id {
			fmt.Println("Running job by name", id)
			s.Jobs[i].Force = true
			if err := j.job.RunNow(); err != nil {
				fmt.Println("Error running job manually", err)
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
			fmt.Println("Deleting forced inactive job", id)
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
			fmt.Println("Deleting forced inactive job", id)
			s.Jobs[i].Running = true
			break
		}
	}
}

func (s *Scheduler) RecreateCtx(name string) {
	for i, j := range s.Jobs {
		if j.job.Name() == name {
			fmt.Println("Recreating context for job", name)
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
	fmt.Println("Rescheduling backups")

	config := s.settings.Config

	for i := range config.Schedules {
		schedule := config.Schedules[i]

		jobDef := gocron.OneTimeJob(gocron.OneTimeJobStartImmediately())

		if schedule.Cron != "" {
			jobDef = gocron.CronJob(schedule.Cron, false)
		}

		j, err := s.gocron.NewJob(
			jobDef,
			gocron.NewTask(func() {

				existing := s.FindJobById(schedule.Id)
				isForced := existing != nil && existing.Force == true
				if !schedule.Active && !isForced {
					fmt.Println("MISSING", schedule.Id)
					return
				}
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
					fmt.Println("before job run")
					s.SetRunningJob(jobName)
				}),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {
						fmt.Println("after job run")
						s.DeleteRunningJob(jobName)
						s.RecreateCtx(jobName)
					},
				),
				gocron.AfterJobRunsWithError(
					func(jobID uuid.UUID, jobName string, err error) {
						fmt.Println("after job run with error", err)
						s.DeleteRunningJob(jobName)
						s.RecreateCtx(jobName)
					},
				),
			))

		if err != nil {
			fmt.Println("Error creating Job", err)
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
