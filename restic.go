package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/go-errors/errors"
	"github.com/labstack/gommon/log"
)

type Restic struct {
	errb *bytes.Buffer
	outb *bytes.Buffer
}

func NewRestic(outb *bytes.Buffer, errb *bytes.Buffer) *Restic {
	r := &Restic{}
	r.errb = errb
	r.outb = outb
	return r
}

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

func (r *Restic) core(repository Repository, cmd ...string) (string, error) {

	cmds := []string{"-r", repository.Path, "--json"}
	cmds = append(cmds, cmd...)
	c := exec.Command("/usr/bin/restic", cmds...)
	c.Stderr = r.errb
	c.Stdout = r.outb
	c.Env = append(c.Env, "RESTIC_PASSWORD="+repository.Password)
	log.Print(c.Env)

	if out, err := c.Output(); err != nil {
		fmt.Println(r.errb.String())
		return "", errors.New(r.errb.String())
	} else {
		fmt.Println(string(out))
		return string(out), nil
	}

}

func (r *Restic) Check(repository Repository) error {
	if _, err := r.core(repository, "check"); err != nil {
		return err
	}
	return nil
}

func (r *Restic) Initialize(repository Repository) error {
	if _, err := r.core(repository, "init"); err != nil {
		return err
	}
	return nil
}

func (r *Restic) Snapshots(repository Repository) []Snapshot {
	if res, err := r.core(repository, "snapshots"); err == nil {
		var data []Snapshot
		if err := json.Unmarshal([]byte(res), &data); err == nil {
			return data
		}
	}

	return []Snapshot{}
}

func (r *Restic) RunBackup(backup Backup, toRepository Repository, fromRepository Repository) {
	fmt.Println("RUNNING BACKUP")
	time.Sleep(30 * time.Second)
	// r.Snapshots(toRepository)
}
