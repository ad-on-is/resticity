package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type Restic struct {
	settings *Settings
	OutputCh *chan ChanMsg
	ErrorCh  *chan ChanMsg
}

func NewRestic(settings *Settings, outch *chan ChanMsg, errch *chan ChanMsg) *Restic {
	r := &Restic{}

	r.settings = settings
	r.OutputCh = outch
	r.ErrorCh = errch
	return r
}

func (r *Restic) PipeOutErr(
	c *exec.Cmd,
	sout *bytes.Buffer,
	serr *bytes.Buffer,
	job *Job,
) {
	stdout, err := c.StdoutPipe()
	if err == nil {
		go func() {
			scanner := bufio.NewScanner(stdout)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				go func(t string) {
					msg := ChanMsg{Id: "", Msg: t, Time: time.Now()}
					if job != nil {
						msg.Id = job.Id
					}
					(*r.OutputCh) <- msg
				}(scanner.Text())

				sout.WriteString(scanner.Text())
			}
		}()

	}

	stderr, err := c.StderrPipe()

	if err == nil {

		go func() {
			scanner := bufio.NewScanner(stderr)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {

				go func(t string) {
					msg := ChanMsg{Id: "", Msg: t, Time: time.Now()}
					if job != nil {
						msg.Id = job.Id
					}
					(*r.OutputCh) <- msg
				}(scanner.Text())
				serr.WriteString(scanner.Text())
			}
		}()
	}
}

func (r *Restic) getEnvs(repository Repository, envs []string) []string {
	envs = append(
		envs,
		[]string{"RESTIC_PASSWORD=" + repository.Password, "RESTIC_PROGRESS_FPS=5"}...)
	if repository.Type == "s3" {
		envs = append(
			envs,
			[]string{
				"AWS_ACCESS_KEY_ID=" + repository.Options.S3Key,
				"AWS_SECRET_ACCESS_KEY=" + repository.Options.S3Secret,
			}...)
	}
	if repository.Type == "azure" {
		envs = append(
			envs,
			[]string{
				"AZURE_ACCOUNT_NAME=" + repository.Options.AzureAccountName,
				"AZURE_ACCOUNT_KEY=" + repository.Options.AzureAccountKey,
				"AZURE_ACCOUNT_SAS=" + repository.Options.AzureAccountSas,
			}...)
	}

	if repository.Type == "gcs" {
		envs = append(
			envs,
			[]string{
				"GOOGLE_PROJECT_ID=" + repository.Options.GoogleProjectId,
				"GOOGLE_APPLICATION_CREDENTIALS=" + repository.Options.GoogleApplicationCredentials,
			}...)
	}
	return envs

}

func (r *Restic) core(
	repository Repository,
	cmd []string,
	envs []string,
	job *Job,
) (string, error) {

	// trigger start

	cmds := []string{"-r", repository.Path, "--json"}
	cmds = append(cmds, cmd...)
	var sout bytes.Buffer
	var serr bytes.Buffer
	var c *exec.Cmd

	if job != nil && job.Ctx != nil {
		c = exec.CommandContext(job.Ctx, "/usr/bin/restic", cmds...)
		if job.Cancel != nil {

			defer (job.Cancel)()
		}
	} else {
		c = exec.Command("/usr/bin/restic", cmds...)
	}

	r.PipeOutErr(c, &sout, &serr, job)

	envs = r.getEnvs(repository, envs)
	log.Info("core", "repo", repository.Path, "cmd", cmd)

	c.Env = append(
		os.Environ(),
		envs...,
	)

	err := c.Start()
	if err != nil {
		log.Error("executing restic command", "err", err)
	}
	c.Wait()

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
			log.Error("browse snapshot: unmarshal", "err", err)
			return []FileDescriptor{}, err
		}
	} else {
		log.Error("browse snapshot: core", "err", err)
		return []FileDescriptor{}, err
	}

}

func (r *Restic) RunSchedule(
	job *Job,
) {

	if job == nil {
		return
	}
	(*r.OutputCh) <- ChanMsg{Id: job.Schedule.Id, Msg: "{\"running\": true}", Time: time.Now()}
	toRepository := r.settings.GetRepositoryById(job.Schedule.ToRepositoryId)
	fromRepository := r.settings.GetRepositoryById(job.Schedule.FromRepositoryId)
	backup := r.settings.GetBackupById(job.Schedule.BackupId)

	switch job.Schedule.Action {
	case "backup":
		if backup == nil || toRepository == nil {
			log.Error("backup", "err", "missing backup and toRepository")
			return
		}
		cmds := []string{"backup", backup.Path, "--tag", "resticity"}
		for _, p := range backup.BackupParams {
			cmds = append(cmds, p...)
		}

		_, err := r.core(*toRepository, cmds, []string{}, job)
		if err != nil {
			log.Error("runschedule", "err", err)
		}
		break
	case "copy-snapshots":
		if fromRepository == nil || toRepository == nil {
			log.Error("copy snapshots", "err", "missing fromRepository and toRepository")
			return
		}
		cmds := []string{"copy"}
		envs := []string{
			"RESTIC_FROM_PASSWORD=" + fromRepository.Password,
			"RESTIC_FROM_REPOSITORY=" + fromRepository.Path,
		}

		r.core(*toRepository, cmds, envs, job)
		break
	case "prune-repository":
		if toRepository == nil {
			log.Error("prune-repository", "err", "missing toRepository")
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
			job,
		)
		if err == nil {
			_, err := r.core(*toRepository, cmds, []string{}, job)
			if err != nil {
				log.Error("prune-repository", "err", err)
			}
		}

		break
	}

}
