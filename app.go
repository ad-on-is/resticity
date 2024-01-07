package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/energye/systray"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	scheduler *Scheduler
	restic    *Restic
	settings  *Settings
}

type BackupJob struct {
	JobId    uuid.UUID `json:"job_id"`
	Schedule Schedule  `json:"schedule"`
}

// NewApp creates a new App application struct
func NewApp(restic *Restic, scheduler *Scheduler, settings *Settings) *App {
	return &App{restic: restic, scheduler: scheduler, settings: settings}
}

func (a *App) toggleSysTrayIcon() {
	default_icon, _ := os.ReadFile(
		"/home/adonis/Development/Go/src/github.com/ad-on-is/resticity/build/appicon.png",
	)
	active_icon, _ := os.ReadFile(
		"/home/adonis/Development/Go/src/github.com/ad-on-is/resticity/build/appicon_active.png",
	)
	def := true
	_, err := a.scheduler.gocron.NewJob(
		gocron.DurationJob(500*time.Millisecond),
		gocron.NewTask(func() {
			fmt.Println("toggle")
			if def {
				if len(a.scheduler.RunningJobs) > 0 {
					systray.SetIcon(active_icon)
				}
				def = false
			} else {
				systray.SetIcon(default_icon)
				def = true
			}
		}),
	)
	if err != nil {
		fmt.Println("Error creating job", err)
	}

}

func (a *App) systemTray() {
	ico, _ := os.ReadFile(
		"/home/adonis/Development/Go/src/github.com/ad-on-is/resticity/build/appicon.png",
	)

	systray.CreateMenu()

	systray.SetIcon(ico) // read the icon from a file

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
	a.toggleSysTrayIcon()
	go systray.Run(a.systemTray, func() {})

}

func (a *App) GetBackupJobs() []BackupJob {
	return a.scheduler.RunningJobs
}

func (a *App) Snapshots(id string) []Snapshot {
	s := a.settings.Config
	var r *Repository
	for i := range s.Repositories {
		if s.Repositories[i].Id == id {
			r = &s.Repositories[i]
		}
	}
	fmt.Println(r)
	if r != nil {
		return a.restic.Snapshots(*r)
	}
	return []Snapshot{}

}

func (a *App) StopBackup(id uuid.UUID) {
	// a.scheduler.RemoveJob(id)
	// a.RescheduleBackups()
}

func (a *App) CheckRepository(r Repository) string {
	files, err := os.ReadDir(r.Path)
	if err != nil {
		return err.Error()
	}
	if len(files) > 0 {
		if err := a.restic.Check(r); err != nil {
			return err.Error()
		} else {
			return "REPO_OK_EXISTING"
		}
	}

	return "REPO_OK_EMPTY"
}

func (a *App) InitializeRepository(r Repository) string {
	if err := a.restic.Initialize(r); err != nil {
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
