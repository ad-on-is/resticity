package restic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/go-errors/errors"
	"github.com/labstack/gommon/log"
)

type B2Options struct {
	B2AccountId  string `json:"b2_account_id"`
	B2AccountKey string `json:"b2_account_key"`
}

type AzureOptions struct {
	AzureAccountName    string `json:"azure_account_name"`
	AzureAccountKey     string `json:"azure_account_key"`
	AzureAccountSas     string `json:"azure_account_sas"`
	AzureEndpointSuffix string `json:"azure_endpoint_suffix"`
}

type Options struct {
	B2Options
	AzureOptions
}

type RepositoryType int32

const (
	LOCAL  RepositoryType = iota
	B2     RepositoryType = iota
	AWS    RepositoryType = iota
	AZURE  RepositoryType = iota
	GOOGLE RepositoryType = iota
)

type Param struct {
	k string
	v string
}

type BackupParam struct {
	Param
}

type PruneParam struct {
	Param
}

type Repository struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Type        RepositoryType `json:"type"`
	PruneParams []Param        `json:"prune_params"`
	Path        string         `json:"path"`
	Password    string         `json:"password"`
	Options     Options        `json:"options"`
}

type Snapshot struct {
	Id             string   `json:"id"`
	Time           string   `json:"time"`
	Paths          []string `json:"paths"`
	Hostname       string   `json:"hostname"`
	Username       string   `json:"username"`
	UID            uint32   `json:"uid"`
	GID            uint32   `json:"gid"`
	ShortId        string   `json:"short_id"`
	Tags           []string `json:"tags"`
	ProgramVersion string   `json:"program_version"`
}

type Backup struct {
	Id           string     `json:"id"`
	Path         string     `json:"path"`
	Name         string     `json:"name"`
	Cron         string     `json:"cron"`
	BackupParams [][]string `json:"backup_params"`
	Targets      []string   `json:"targets"`
}

func core(r Repository, cmd ...string) (string, error) {
	var errb bytes.Buffer
	cmds := []string{"-r", r.Path, "--json"}
	cmds = append(cmds, cmd...)
	c := exec.Command("/usr/bin/restic", cmds...)
	c.Stderr = &errb
	c.Env = append(c.Env, "RESTIC_PASSWORD="+r.Password)
	log.Print(c.Env)

	if out, err := c.Output(); err != nil {
		fmt.Println(errb.String())
		return "", errors.New(errb.String())
	} else {
		fmt.Println(string(out))
		return string(out), nil
	}

}

func Check(r Repository) error {
	if _, err := core(r, "check"); err != nil {
		return err
	}
	return nil
}

func Initialize(r Repository) error {
	if _, err := core(r, "init"); err != nil {
		return err
	}
	return nil
}

func Snapshots(r Repository) []Snapshot {
	if res, err := core(r, "snapshots"); err == nil {
		var data []Snapshot
		if err := json.Unmarshal([]byte(res), &data); err == nil {
			return data
		}
	}

	return []Snapshot{}
}
