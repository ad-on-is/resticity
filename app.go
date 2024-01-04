package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"resticity/restic"
	"sync"
	"time"

	"github.com/adrg/xdg"
	"github.com/energye/systray"
	"github.com/energye/systray/icon"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	jmu         sync.Mutex
	scheduler   gocron.Scheduler
	runningJobs []BackupJob
}

type BackupJob struct {
	BackupId     string    `json:"backup_id"`
	RepositoryId string    `json:"repository_id"`
	JobId        uuid.UUID `json:"job_id"`
}

type Settings struct {
	Repositories []restic.Repository `json:"repositories"`
	Backups      []restic.Backup     `json:"backups"`
	Schedules    []restic.Schedule   `json:"schedules"`
}

func settingsFile() string {
	return xdg.ConfigHome + "/resticity.json"
}

func GetSettings() Settings {
	data := Settings{}
	if file, err := os.Open(settingsFile()); err == nil {
		if str, err := io.ReadAll(file); err == nil {
			if err := json.Unmarshal([]byte(str), &data); err == nil {
				return data
			}
		}
	} else {
		fmt.Println("error", err)
	}
	return data
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) systemTray() {
	ico, _ := os.ReadFile("/home/adonis/Development/Flutter/Sdk/dev/docs/favicon.ico")

	systray.SetIcon(icon.Data) // read the icon from a file
	fmt.Println(len(ico))

	systray.SetTitle("resticity")
	systray.SetTooltip("Resticity")

	show := systray.AddMenuItem("Open resticity", "Show the main window")
	systray.AddSeparator()

	exit := systray.AddMenuItem("Quit", "Quit resticity")

	show.Click(func() {

		runtime.WindowShow(a.ctx)
	})
	exit.Click(func() { os.Exit(0) })

	systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(a.ctx) })
	// systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
	systray.SetOnRClick(func(menu systray.IMenu) { runtime.WindowHide(a.ctx) })
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go systray.Run(a.systemTray, func() {})
	if s, err := gocron.NewScheduler(); err == nil {
		go s.Start()
		a.scheduler = s
		a.RescheduleBackups()
	}

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetBackupJobs() []BackupJob {
	return a.runningJobs
}

func (a *App) Snapshots(id string) []restic.Snapshot {
	s := GetSettings()
	var r *restic.Repository
	for i := range s.Repositories {
		if s.Repositories[i].Id == id {
			r = &s.Repositories[i]
		}
	}
	fmt.Println(r)
	if r != nil {
		return restic.Snapshots(*r)
	}
	return []restic.Snapshot{}

}

func (a *App) Settings() Settings {
	return GetSettings()
}

func (a *App) SaveSettings(data Settings) {
	a.RescheduleBackups()
	fmt.Println("Saving settings")
	if str, err := json.MarshalIndent(data, " ", " "); err == nil {
		fmt.Println("Settings saved")
		if err := os.WriteFile(settingsFile(), str, 0644); err != nil {
			fmt.Println("error", err)
		} else {
			a.RescheduleBackups()
		}
	} else {
		fmt.Println("error", err)
	}
}

func (a *App) StopBackup(id uuid.UUID) {
	a.scheduler.RemoveJob(id)
	a.RescheduleBackups()
}

func (a *App) RescheduleBackups() {
	s := GetSettings()
	for _, b := range s.Backups {

		// todo: run missed backups
		/*
			- get last snapshot
			- parse cron as date/duration whatever
			- compare last snapshot date with cron date and job.NextRun()
		*/

		for _, t := range b.Targets {

			jobName := "BACKUP-" + b.Id + "-TARGET-" + t
			var job gocron.Job
			for j := range a.scheduler.Jobs() {
				if a.scheduler.Jobs()[j].Name() == jobName {
					job = a.scheduler.Jobs()[j]
					break
				}
			}
			if job != nil {
				a.scheduler.RemoveJob(job.ID())
			}
			if b.Cron == "" {
				continue
			}
			a.scheduler.NewJob(
				gocron.CronJob(b.Cron, false),
				gocron.NewTask(func(backup restic.Backup) {
					// actual backup
					log.Print("doing job", backup.Name)
					time.Sleep(30 * time.Second)
				}, b),
				gocron.WithTags("backup:"+b.Id, "repository:"+t),
				gocron.WithName(jobName),
				gocron.WithEventListeners(
					gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
						a.jmu.Lock()
						defer a.jmu.Unlock()
						a.runningJobs = append(
							a.runningJobs,
							BackupJob{
								BackupId:     b.Id,
								RepositoryId: t,
								JobId:        jobID,
							},
						)
					}),
					gocron.AfterJobRuns(
						func(jobID uuid.UUID, jobName string) {
							a.jmu.Lock()
							defer a.jmu.Unlock()
							for i := range a.runningJobs {
								if a.runningJobs[i].BackupId == b.Id &&
									a.runningJobs[i].RepositoryId == t {
									a.runningJobs = append(
										a.runningJobs[:i],
										a.runningJobs[i+1:]...)
									break
								}
							}
							// do something after the job completes

						},
					),
					gocron.AfterJobRunsWithError(
						func(jobID uuid.UUID, jobName string, err error) {
							a.jmu.Lock()
							defer a.jmu.Unlock()
							for i := range a.runningJobs {
								if a.runningJobs[i].BackupId == b.Id &&
									a.runningJobs[i].RepositoryId == t {
									a.runningJobs = append(
										a.runningJobs[:i],
										a.runningJobs[i+1:]...)
									break
								}
							}
						},
					),
				),
			)
		}

	}

}

func (a *App) CheckRepository(r restic.Repository) string {
	files, err := os.ReadDir(r.Path)
	if err != nil {
		return err.Error()
	}
	if len(files) > 0 {
		if err := restic.Check(r); err != nil {
			return err.Error()
		} else {
			return "REPO_OK_EXISTING"
		}
	}

	return "REPO_OK_EMPTY"
}

func (a *App) InitializeRepository(r restic.Repository) string {
	if err := restic.Initialize(r); err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) SelectDirectory(title string) string {
	if dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	}); err == nil {
		return dir
	}

	return ""
}
