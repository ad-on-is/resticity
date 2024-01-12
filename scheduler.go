package main

import (
	"fmt"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler struct {
	gocron      gocron.Scheduler
	restic      *Restic
	RunningJobs []BackupJob
	Jobs        []gocron.Job
	ManualJobs  []string
	jmu         sync.Mutex
	settings    *Settings
}

func NewScheduler(settings *Settings, restic *Restic) (*Scheduler, error) {

	s := &Scheduler{}
	s.settings = settings
	s.restic = restic
	s.ManualJobs = []string{}
	s.RunningJobs = []BackupJob{}
	if gc, err := gocron.NewScheduler(); err == nil {
		s.gocron = gc
		s.gocron.Start()
		return s, nil
	} else {
		return nil, err
	}

}

func (s *Scheduler) RunJobByName(name string) {
	for _, job := range s.Jobs {
		if job.Name() == name {
			s.ManualJobs = append(s.ManualJobs, name)
			if err := job.RunNow(); err != nil {
				fmt.Println("Error running job manually", err)
			}
			break
		}
	}
}

func (s *Scheduler) DeleteBackgroundJob(jobID uuid.UUID) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.RunningJobs {
		if j.JobId == jobID {
			s.RunningJobs = append(
				s.RunningJobs[:i],
				s.RunningJobs[i+1:]...)
			break
		}
	}
}

func (s *Scheduler) DeleteManualJob(name string) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i, j := range s.ManualJobs {
		if j == name {
			s.ManualJobs = append(
				s.ManualJobs[:i],
				s.ManualJobs[i+1:]...)
			break
		}
	}
}

func (s *Scheduler) GetRunningJobs() []BackupJob {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	return s.RunningJobs
}

func (s *Scheduler) RescheduleBackups() {

	s.Jobs = []gocron.Job{}

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
				toRepository := s.settings.GetRepositoryById(schedule.ToRepositoryId)
				fromRepository := s.settings.GetRepositoryById(schedule.FromRepositoryId)
				backup := s.settings.GetBackupById(schedule.BackupId)
				if !schedule.Active && !StringArrayContains(s.ManualJobs, schedule.Id) {
					return
				}
				s.restic.RunBackup(schedule.Action, backup, toRepository, fromRepository)
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
						BackupJob{
							JobId:    jobID,
							Schedule: schedule,
						},
					)
				}),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {
						s.DeleteBackgroundJob(jobID)
						s.DeleteManualJob(jobName)
					},
				),
				gocron.AfterJobRunsWithError(
					func(jobID uuid.UUID, jobName string, err error) {
						s.DeleteBackgroundJob(jobID)
						s.DeleteManualJob(jobName)
					},
				),
			))

		if err != nil {
			fmt.Println("Error creating Job", err)
			continue
		}

		s.Jobs = append(s.Jobs, j)

	}

}
