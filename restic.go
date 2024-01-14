package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Restic struct {
	errb     *bytes.Buffer
	outb     *bytes.Buffer
	settings *Settings
}

func NewRestic(errb *bytes.Buffer, outb *bytes.Buffer, settings *Settings) *Restic {
	r := &Restic{}
	r.errb = errb
	r.outb = outb
	r.settings = settings
	return r
}

func (r *Restic) core(
	repository Repository,
	cmd []string,
	envs []string,
	ctx *context.Context,
	cancel *context.CancelFunc,
	ch *chan string,
) (string, error) {

	cmds := []string{"-r", repository.Path, "--json"}
	cmds = append(cmds, cmd...)
	var sout bytes.Buffer
	var serr bytes.Buffer
	var c *exec.Cmd
	if ctx != nil {
		c = exec.CommandContext(*ctx, "/usr/bin/restic", cmds...)
		if cancel != nil {

			defer (*cancel)()
		}
	} else {
		c = exec.Command("/usr/bin/restic", cmds...)
	}

	stdout, err := c.StdoutPipe()

	if err == nil {
		go func() {
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {

				if ch != nil {
					// fmt.Println(t)
					go func(t string) {
						*ch <- t
					}(scanner.Text())
				}
				sout.WriteString(scanner.Text())
			}
		}()
	}

	c.Env = append(
		os.Environ(),
		"RESTIC_PASSWORD="+repository.Password,
		"RESTIC_PROGRESS_FPS=5",
	)

	err = c.Start()
	if err != nil {
		fmt.Println(err)
	}
	c.Wait()
	go func() {
		if ch != nil {
			*ch <- ""
		}
	}()
	// fmt.Println(sp.Data)
	r.errb.Write(serr.Bytes())
	r.outb.Write(sout.Bytes())

	if serr.Len() > 0 {
		return "", errors.New(serr.String())
	}

	return sout.String(), nil

}

func (r *Restic) Exec(repository Repository, cmds []string, envs []string) (string, error) {
	if data, err := r.core(repository, cmds, envs, nil, nil, nil); err != nil {
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
	if res, err := r.core(repository, []string{"ls", "-l", "--human-readable", snapshotId, path}, []string{}, nil, nil, nil); err == nil {
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
	job *Job,
) {

	if job == nil {
		return
	}
	toRepository := r.settings.GetRepositoryById(job.Schedule.ToRepositoryId)
	fromRepository := r.settings.GetRepositoryById(job.Schedule.FromRepositoryId)
	backup := r.settings.GetBackupById(job.Schedule.BackupId)

	switch job.Schedule.Action {
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

		_, err := r.core(*toRepository, cmds, []string{}, &job.Ctx, &job.Cancel, &job.Chan)
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
		_, err := r.core(
			*toRepository,
			[]string{"unlock"},
			[]string{},
			&job.Ctx,
			&job.Cancel,
			&job.Chan,
		)
		if err == nil {
			_, err := r.core(*toRepository, cmds, []string{}, &job.Ctx, &job.Cancel, &job.Chan)
			if err != nil {
				fmt.Println(err)
			}
		}

		break
	}

}
