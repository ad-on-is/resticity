package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/ad-on-is/resticity/internal"
	"github.com/adrg/xdg"

	"github.com/charmbracelet/log"
	"github.com/energye/systray"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	scheduler *internal.Scheduler
	restic    *internal.Restic
	settings  *internal.Settings
	assets    *embed.FS
}

// NewApp creates a new App application struct
func NewApp(
	restic *internal.Restic,
	scheduler *internal.Scheduler,
	settings *internal.Settings,
	assets *embed.FS,
) *App {
	return &App{restic: restic, scheduler: scheduler, settings: settings}
}

func (a *App) SaveIcon(icon []byte, file string) {
	path := filepath.Join(xdg.CacheHome, "resticity")
	os.WriteFile(filepath.Join(path, file), icon, 0644)
}

func (a *App) toggleSysTrayIcon() {
	default_icon, err := assets.ReadFile(
		"frontend/.output/public/appicon.png",
	)
	if err != nil {
		log.Error(err)
	} else {
		a.SaveIcon(default_icon, "appicon.png")

	}
	active_icon, err := assets.ReadFile(
		"frontend/.output/public/appicon_active.png",
	)

	if err != nil {
		log.Error(err)
	} else {
		a.SaveIcon(active_icon, "appicon_active.png")

	}

	_, err = a.scheduler.Gocron.NewJob(
		gocron.DurationJob(500*time.Millisecond),
		gocron.NewTask(func() {
			running := funk.Filter(
				a.scheduler.Jobs,
				func(j internal.Job) bool { return j.Running == true },
			).([]internal.Job)
			if len(running) > 0 {
				systray.SetIcon(active_icon)
			} else {
				systray.SetIcon(default_icon)
			}

		}),
	)
	if err != nil {
		log.Error("Error creating job", err)
	}

}

func (a *App) systemTray() {
	ico, _ := os.ReadFile(
		"/home/adonis/Development/Go/resticity/build/appicon.png",
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

func (a *App) StopBackup(id uuid.UUID) {
	// a.scheduler.RemoveJob(id)
	// a.RescheduleBackups()
}

func (a *App) SelectDirectory(title string) string {
	if dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	}); err == nil {
		return dir
	}

	return ""
}

func (a *App) SelectFile(title string) string {
	if dir, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	}); err == nil {
		return dir
	}

	return ""
}

func (a *App) FakeCreateForModels() (internal.SnapshotGroup, internal.Repository, internal.Backup, internal.Config, internal.Schedule, internal.FileDescriptor, internal.ScheduleObject) {
	return internal.SnapshotGroup{}, internal.Repository{}, internal.Backup{}, internal.Config{}, internal.Schedule{}, internal.FileDescriptor{}, internal.ScheduleObject{}
}
