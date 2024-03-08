package main

import (
	"context"
	"embed"
	"os"
	"sync"
	"time"

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
	scheduler *Scheduler
	restic    *Restic
	settings  *Settings
	assets    *embed.FS
}

type Settings struct {
	file   string
	Config Config `json:"config"`
	mux    sync.Mutex
}

type S3Options struct {
	S3Key    string `json:"s3_key"`
	S3Secret string `json:"s3_secret"`
}

type AzureOptions struct {
	AzureAccountName string `json:"azure_account_name"`
	AzureAccountKey  string `json:"azure_account_key"`
	AzureAccountSas  string `json:"azure_account_sas"`
}

type GcsOptions struct {
	GoogleProjectId              string `json:"google_project_id"`
	GoogleApplicationCredentials string `json:"google_application_credentials"`
}

type Options struct {
	S3Options
	AzureOptions
	GcsOptions
}

type GroupKey struct {
	Hostname string   `json:"hostname"`
	Paths    []string `json:"paths"`
	Tags     []string `json:"tags"`
}

type SnapshotGroup struct {
	GroupKey  GroupKey   `json:"group_key"`
	Snapshots []Snapshot `json:"snapshots"`
}

type Snapshot struct {
	Id             string    `json:"id"`
	Time           time.Time `json:"time"`
	Paths          []string  `json:"paths"`
	Hostname       string    `json:"hostname"`
	Username       string    `json:"username"`
	UID            uint32    `json:"uid"`
	GID            uint32    `json:"gid"`
	ShortId        string    `json:"short_id"`
	Tags           []string  `json:"tags"`
	ProgramVersion string    `json:"program_version"`
}

type FileDescriptor struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Path  string `json:"path"`
	Size  uint32 `json:"size"`
	Mtime string `json:"mtime"`
}

type Mount struct {
	RepositoryId string `json:"repository_id"`
	Path         string `json:"path"`
}

type Repository struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	PruneParams [][]string `json:"prune_params"`
	Path        string     `json:"path"`
	Password    string     `json:"password"`
	Options     Options    `json:"options"`
}

type Backup struct {
	Id           string     `json:"id"`
	Path         string     `json:"path"`
	Name         string     `json:"name"`
	Cron         string     `json:"cron"`
	BackupParams [][]string `json:"backup_params"`
	Targets      []string   `json:"targets"`
}

type Schedule struct {
	Id               string `json:"id"`
	Action           string `json:"action"`
	BackupId         string `json:"backup_id"`
	ToRepositoryId   string `json:"to_repository_id"`
	FromRepositoryId string `json:"from_repository_id"`
	Cron             string `json:"cron"`
	Active           bool   `json:"active"`
	LastRun          string `json:"last_run"`
	LastError        string `json:"last_error"`
}

type Config struct {
	Mounts       []Mount      `json:"mounts"`
	Repositories []Repository `json:"repositories"`
	Backups      []Backup     `json:"backups"`
	Schedules    []Schedule   `json:"schedules"`
}

// NewApp creates a new App application struct
func NewApp(restic *Restic, scheduler *Scheduler, settings *Settings, assets *embed.FS) *App {
	return &App{restic: restic, scheduler: scheduler, settings: settings}
}

func (a *App) toggleSysTrayIcon() {
	default_icon, _ := os.ReadFile(
		"/home/adonis/Development/Go/resticity/frontend/public/appicon.png",
	)
	active_icon, _ := os.ReadFile(
		"/home/adonis/Development/Go/resticity/frontend/public/appicon_active.png",
	)

	_, err := a.scheduler.gocron.NewJob(
		gocron.DurationJob(500*time.Millisecond),
		gocron.NewTask(func() {
			running := funk.Filter(
				a.scheduler.Jobs,
				func(j Job) bool { return j.Running == true },
			).([]Job)
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

func (a *App) FakeCreateForModels() (SnapshotGroup, Repository, Backup, Config, Schedule, FileDescriptor) {
	return SnapshotGroup{}, Repository{}, Backup{}, Config{}, Schedule{}, FileDescriptor{}
}
