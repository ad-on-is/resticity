package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"resticity/restic"

	"github.com/adrg/xdg"
	"github.com/energye/systray"
	"github.com/energye/systray/icon"
	"github.com/go-co-op/gocron/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	scheduler gocron.Scheduler
}

type Settings struct {
	Repositories []restic.Repository `json:"repositories"`
	Backups      []restic.Backup     `json:"backups"`
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
	systray.SetTooltip("yeaaah baby")

	show := systray.AddMenuItem("Show", "Show The Window")
	systray.AddSeparator()
	// exit := systray.AddMenuItem("Exit", "Quit The Program")cd

	show.Click(func() {

		runtime.WindowToggleMaximise(a.ctx)
	})
	// exit.Click(func() { os.Exit(0) })

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

func (a *App) Snapshots(id string) []restic.Snapshot {
	fmt.Println("IIIID", id)
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

func (a *App) RescheduleBackups() {
	s := GetSettings()
	for b := range s.Backups {
		var job gocron.Job
		for j := range a.scheduler.Jobs() {
			if a.scheduler.Jobs()[j].Name() == s.Backups[b].Name {
				job = a.scheduler.Jobs()[j]
				break
			}
		}
		if job != nil {
			a.scheduler.RemoveJob(job.ID())
		}
		// todo: run missed backups
		/*
			- get last snapshot
			- parse cron as date/duration whatever
			- compare last snapshot date with cron date and job.NextRun()
		*/

		a.scheduler.NewJob(
			gocron.CronJob("*/1 * * * *", false),
			gocron.NewTask(func(backup restic.Backup) {
				log.Print("doing job", backup.Name)
			}, s.Backups[b]),
			gocron.WithTags("backup", "bdonis"),
			gocron.WithName(s.Backups[b].Name),
		)
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
