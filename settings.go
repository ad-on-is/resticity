package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/adrg/xdg"
)

type Settings struct {
	Config Config
}

type Repository struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Type        RepositoryType `json:"type"`
	PruneParams [][]string     `json:"prune_params"`
	Path        string         `json:"path"`
	Password    string         `json:"password"`
	Options     Options        `json:"options"`
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
	BackupId         string `json:"backup_id"`
	ToRepositoryId   string `json:"to_repository_id"`
	FromRepositoryId string `json:"from_repository_id"`
	Cron             string `json:"cron"`
	Active           bool   `json:"active"`
}

type Config struct {
	Repositories []Repository `json:"repositories"`
	Backups      []Backup     `json:"backups"`
	Schedules    []Schedule   `json:"schedules"`
}

func NewSettings() *Settings {
	s := &Settings{}
	s.Config = s.readFile()
	return s
}

func (s *Settings) GetRepositoryById(id string) *Repository {
	for _, r := range s.Config.Repositories {
		if r.Id == id {
			return &r
		}
	}
	return nil
}

func (s *Settings) GetBackupById(id string) *Backup {
	for _, b := range s.Config.Backups {
		if b.Id == id {
			return &b
		}
	}
	return nil
}

func (s *Settings) readFile() Config {
	data := Config{}
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

func settingsFile() string {
	return xdg.ConfigHome + "/resticity.json"
}

func (s *Settings) Save(data Config) error {
	s.Config = data
	fmt.Println("Saving settings")
	if str, err := json.MarshalIndent(s.Config, " ", " "); err == nil {
		fmt.Println("Settings saved")
		if err := os.WriteFile(settingsFile(), str, 0644); err != nil {
			fmt.Println("error", err)
			return err
		}
	} else {
		fmt.Println("error", err)
		return err
	}
	return nil
}
