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
	jmu         sync.Mutex
	settings    *Settings
}

func NewScheduler(settings *Settings, restic *Restic) (*Scheduler, error) {

	s := &Scheduler{}
	s.settings = settings
	s.restic = restic
	if scheduler, err := gocron.NewScheduler(); err == nil {
		s.gocron = scheduler
		s.gocron.Start()
		return s, nil
	} else {
		return nil, err
	}

}

func (s *Scheduler) DeleteBackgroundJob(jobID uuid.UUID) {
	s.jmu.Lock()
	defer s.jmu.Unlock()
	for i := range s.RunningJobs {
		if s.RunningJobs[i].JobId == jobID {
			s.RunningJobs = append(
				s.RunningJobs[:i],
				s.RunningJobs[i+1:]...)
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
	fmt.Println("Rescheduling backups")
	config := s.settings.Config

	for i := range config.Schedules {
		schedule := config.Schedules[i]
		fmt.Println("SCHEDULE", s)
		var job gocron.Job
		for j := range s.gocron.Jobs() {
			if s.gocron.Jobs()[j].Name() == schedule.Id {
				job = s.gocron.Jobs()[j]
			}
		}
		if job != nil {
			fmt.Println("Found job")
			s.gocron.RemoveJob(job.ID())
		}
		if schedule.Cron == "" {
			continue
		}

		_, err := s.gocron.NewJob(
			gocron.CronJob(schedule.Cron, false),
			gocron.NewTask(func() {
				var toRepository Repository
				var fromRepository Repository
				var backup Backup
				for _, r := range s.settings.Config.Repositories {
					if r.Id == schedule.ToRepositoryId {
						toRepository = r
					}
					if r.Id == schedule.FromRepositoryId {
						fromRepository = r
					}
				}
				for _, b := range s.settings.Config.Backups {
					if b.Id == schedule.BackupId {
						backup = b
						break
					}

				}
				s.restic.RunBackup(backup, toRepository, fromRepository)
			}),
			gocron.WithName(schedule.Id),
			gocron.WithTags(
				"backup:"+schedule.BackupId,
				"repository:"+schedule.ToRepositoryId,
				"fromrepository:"+schedule.ToRepositoryId,
			),
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
					s.jmu.Lock()
					defer s.jmu.Unlock()
					s.RunningJobs = append(
						s.RunningJobs,
						BackupJob{
							JobId:    jobID,
							Schedule: schedule,
						},
					)
					fmt.Println("BEFORE JOB RUNS", len(s.RunningJobs))
				}),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {
						s.DeleteBackgroundJob(jobID)

					},
				),
				gocron.AfterJobRunsWithError(
					func(jobID uuid.UUID, jobName string, err error) {
						s.DeleteBackgroundJob(jobID)
					},
				),
			))

		if err != nil {
			fmt.Println("ERROR", err)
		}
	}

}
