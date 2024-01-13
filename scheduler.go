package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Job struct {
	job      gocron.Job
	schedule Schedule
}

type ScheduleContext struct {
	Id     string
	Ctx    context.Context
	Cancel context.CancelFunc
}

type RunningJob struct {
	JobId    uuid.UUID `json:"job_id"`
	Schedule Schedule  `json:"schedule"`
}

type Scheduler struct {
	gocron            gocron.Scheduler
	restic            *Restic
	RunningJobs       []RunningJob
	ScheduleContexts  []ScheduleContext
	Jobs              []Job
	ForceInactiveJobs []string
	jmu               sync.Mutex
	settings          *Settings
}

func NewScheduler(settings *Settings, restic *Restic) (*Scheduler, error) {

	s := &Scheduler{}
	s.settings = settings
	s.restic = restic
	s.ForceInactiveJobs = []string{}
	s.RunningJobs = []RunningJob{}
	if gc, err := gocron.NewScheduler(); err == nil {
		s.gocron = gc
		s.gocron.Start()
		return s, nil
	} else {
		return nil, err
	}

}

func (s *Scheduler) RunJobByName(name string) {
	fmt.Println("should run", name)
	for _, j := range s.Jobs {
		if j.job.Name() == name {
			fmt.Println("Running job by name", name)
			s.ForceInactiveJobs = append(s.ForceInactiveJobs, name)
			if err := j.job.RunNow(); err != nil {
				fmt.Println("Error running job manually", err)
			}
			break
		}
	}
}

func (s *Scheduler) StopJobByName(name string) {
	for _, c := range s.ScheduleContexts {
		if c.Id == name {
			c.Cancel()
			break
		}
	}
}

func (s *Scheduler) DeleteRunningJob(jobID uuid.UUID) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.RunningJobs {
		if j.JobId == jobID {
			fmt.Println("Deleting forced inactive job", jobID)
			s.RunningJobs = append(
				s.RunningJobs[:i],
				s.RunningJobs[i+1:]...)
			break
		}
	}
}

func (s *Scheduler) RecreateCtx(name string) {
	for i, c := range s.ScheduleContexts {
		if c.Id == name {
			fmt.Println("Recreating context for job", name)
			ctx, cancel := context.WithCancel(context.Background())
			s.ScheduleContexts[i].Ctx = ctx
			s.ScheduleContexts[i].Cancel = cancel
			break
		}
	}
}

func (s *Scheduler) DeleteForcedInactiveJob(name string) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.ForceInactiveJobs {
		if j == name {
			fmt.Println("Deleting forced inactive job", name)
			s.ForceInactiveJobs = append(
				s.ForceInactiveJobs[:i],
				s.ForceInactiveJobs[i+1:]...)
			break
		}
	}
}

func (s *Scheduler) GetRunningJobs() []RunningJob {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	return s.RunningJobs
}

func (s *Scheduler) GetContextById(id string) *ScheduleContext {
	for _, sc := range s.ScheduleContexts {
		if sc.Id == id {
			return &sc
		}
	}
	return nil
}

func (s *Scheduler) RescheduleBackups() {

	s.Jobs = []Job{}
	s.ScheduleContexts = []ScheduleContext{}
	s.RunningJobs = []RunningJob{}
	s.ForceInactiveJobs = []string{}
	fmt.Println("Rescheduling backups")

	config := s.settings.Config

	for i := range config.Schedules {
		schedule := config.Schedules[i]
		ctx, cancel := context.WithCancel(context.Background())
		sctx := ScheduleContext{Id: schedule.Id, Ctx: ctx, Cancel: cancel}
		s.ScheduleContexts = append(s.ScheduleContexts, sctx)
		jobDef := gocron.OneTimeJob(gocron.OneTimeJobStartImmediately())

		if schedule.Cron != "" {
			jobDef = gocron.CronJob(schedule.Cron, false)
		}

		j, err := s.gocron.NewJob(
			jobDef,
			gocron.NewTask(func() {
				toRepository := s.settings.GetRepositoryById(schedule.ToRepositoryId)
				fromRepository := s.settings.GetRepositoryById(schedule.FromRepositoryId)
				backup := s.settings.GetBackupById(schedule.BackupId)
				if !schedule.Active && !StringArrayContains(s.ForceInactiveJobs, schedule.Id) {
					fmt.Println("MISSING", schedule.Id)
					return
				}
				s.restic.RunSchedule(
					s.GetContextById(schedule.Id),
					schedule.Action,
					backup,
					toRepository,
					fromRepository,
				)
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
					s.jmu.Lock()
					defer s.jmu.Unlock()
					s.RunningJobs = append(
						s.RunningJobs,
						RunningJob{
							JobId:    jobID,
							Schedule: schedule,
						},
					)
				}),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {
						fmt.Println("after job run")
						s.DeleteRunningJob(jobID)
						s.DeleteForcedInactiveJob(jobName)
						s.RecreateCtx(jobName)
					},
				),
				gocron.AfterJobRunsWithError(
					func(jobID uuid.UUID, jobName string, err error) {
						fmt.Println("after job run with error", err)
						s.DeleteRunningJob(jobID)
						s.DeleteForcedInactiveJob(jobName)
						s.RecreateCtx(jobName)
					},
				),
			))

		if err != nil {
			fmt.Println("Error creating Job", err)
			continue
		}

		s.Jobs = append(
			s.Jobs,
			Job{job: j, schedule: schedule},
		)

	}

}
