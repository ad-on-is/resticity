package internal

import (
	"sync"
	"time"
)

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

type Repository struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	PruneParams  [][]string `json:"prune_params"`
	Path         string     `json:"path"`
	Password     string     `json:"password"`
	PasswordFile string     `json:"password_file"`
	Options      Options    `json:"options"`
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

type AppSettingsNotifications struct {
	OnScheduleError   bool `json:"on_schedule_error"`
	OnScheduleSuccess bool `json:"on_schedule_success"`
	OnScheduleStart   bool `json:"on_schedule_start"`
}

type AppSettingsHooks struct {
	OnScheduleError   string `json:"on_schedule_error"`
	OnScheduleSuccess string `json:"on_schedule_success"`
	OnScheduleStart   string `json:"on_schedule_start"`
}

type AppSettings struct {
	Theme                 string                   `json:"theme"`
	PreserveErrorLogsDays uint32                   `json:"preserve_error_logs_days"`
	Hooks                 AppSettingsHooks         `json:"hooks"`
	Notifications         AppSettingsNotifications `json:"notifications"`
}

type Config struct {
	Repositories []Repository `json:"repositories"`
	Backups      []Backup     `json:"backups"`
	Schedules    []Schedule   `json:"schedules"`
	AppSettings  AppSettings  `json:"app_settings"`
}

type BrowseData struct {
	Path string `json:"path"`
}

type MountData struct {
	Path string `json:"path"`
}

type RestoreData struct {
	RootPath string `json:"root_path"`
	FromPath string `json:"from_path"`
	ToPath   string `json:"to_path"`
}

type Output struct {
	Id  string `json:"id"`
	Out any    `json:"out"`
}

type MsgJob struct {
	Id       string   `json:"id"`
	Schedule Schedule `json:"schedule"`
	Running  bool     `json:"running"`
	Force    bool     `json:"force"`
}

type ChanMsg struct {
	Id   string
	Msg  string
	Time time.Time
}

type JobMsg struct {
	Id   string    `json:"id"`
	Out  string    `json:"out"`
	Err  string    `json:"err"`
	Time time.Time `json:"time"`
}

type MountMsg struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type ScheduleObject struct {
	Schedule       Schedule    `json:"schedule"`
	ToRepository   *Repository `json:"to_repository"`
	FromRepository *Repository `json:"from_repository"`
	Backup         *Backup     `json:"backup"`
}
