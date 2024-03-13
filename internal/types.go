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

type WsMsg struct {
	Id   string    `json:"id"`
	Out  string    `json:"out"`
	Err  string    `json:"err"`
	Time time.Time `json:"time"`
}
