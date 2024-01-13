package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Restic struct {
	errb *bytes.Buffer
	outb *bytes.Buffer
}

func NewRestic(errb *bytes.Buffer, outb *bytes.Buffer) *Restic {
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

func (r *Restic) core(
	repository Repository,
	cmd []string,
	envs []string,
	sctx *ScheduleContext,
) (string, error) {

	cmds := []string{"-r", repository.Path, "--json"}
	cmds = append(cmds, cmd...)
	var sout bytes.Buffer
	var serr bytes.Buffer
	var c *exec.Cmd
	if sctx != nil {
		c = exec.CommandContext(sctx.Ctx, "/usr/bin/restic", cmds...)
		defer sctx.Cancel()
	} else {
		c = exec.Command("/usr/bin/restic", cmds...)
	}
	c.Stderr = &serr
	c.Stdout = &sout
	c.Env = append(
		os.Environ(),
		"RESTIC_PASSWORD="+repository.Password,
	)

	err := c.Start()
	if err != nil {
		fmt.Println(err)
	}
	c.Wait()
	r.errb.Write(serr.Bytes())
	r.outb.Write(sout.Bytes())

	if serr.Len() > 0 {
		return "", errors.New(serr.String())
	}

	return sout.String(), nil

}

func (r *Restic) Exec(repository Repository, cmds []string, envs []string) (string, error) {
	if data, err := r.core(repository, cmds, envs, nil); err != nil {
		return "", err
	} else {
		return data, nil
	}
}

func (r *Restic) BrowseSnapshot(
	repository Repository,
	snapshotId string,
	path string,
) ([]FileDescriptor, error) {
	if res, err := r.core(repository, []string{"ls", "-l", "--human-readable", snapshotId, path}, []string{}, nil); err == nil {
		res = strings.ReplaceAll(res, "}", "},")
		res = strings.ReplaceAll(res, "\n", "")
		res = "[" + res + "]"
		res = strings.ReplaceAll(res, ",]", "]")
		var data []FileDescriptor
		if err := json.Unmarshal([]byte(res), &data); err == nil {
			return data, nil
		} else {
			fmt.Println("Error parsing JSON", err)
			return []FileDescriptor{}, err
		}
	} else {
		fmt.Println("Error browsing snapshots", err)
		return []FileDescriptor{}, err
	}

}

func (r *Restic) RunSchedule(
	sctx *ScheduleContext,
	action string,
	backup *Backup,
	toRepository *Repository,
	fromRepository *Repository,
) {

	switch action {
	case "backup":
		if backup == nil || toRepository == nil {
			fmt.Println("Nope!")
			return
		}
		cmds := []string{"backup", backup.Path, "--tag", "resticity"}
		for _, p := range backup.BackupParams {
			cmds = append(cmds, p...)
		}

		fmt.Println(cmds)

		_, err := r.core(*toRepository, cmds, []string{}, sctx)
		if err != nil {
			fmt.Println(err)
		}
		break
	case "copy-snapshots":
		if fromRepository == nil || toRepository == nil {
			fmt.Println("Nope!")
			return
		}
		cmds := []string{"copy", "--from-repo", fromRepository.Path}
		envs := []string{"RESTIC_FROM_PASSWORD", fromRepository.Password}
		fmt.Println(cmds)
		fmt.Println(envs)
		// r.core(*toRepository, cmds, envs)
		break
	case "prune-repository":
		if toRepository == nil {
			fmt.Println("Nope!")
			return
		}
		cmds := []string{"forget", "--prune"}
		for _, p := range toRepository.PruneParams {
			cmds = append(cmds, p...)
		}
		_, err := r.core(*toRepository, []string{"unlock"}, []string{}, sctx)
		if err == nil {
			_, err := r.core(*toRepository, cmds, []string{}, sctx)
			if err != nil {
				fmt.Println(err)
			}
		}

		break
	}

}
