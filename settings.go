package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func NewSettings(flagFile string) *Settings {
	s := &Settings{}
	s.file = flagFile
	if s.file == "" {
		s.file = os.Getenv("RESTICITY_SETTINGS_FILE")
	}
	if s.file == "" {
		s.file = filepath.Join(xdg.ConfigHome, "resticity", "config.json")
	}

	if _, err := os.Stat(s.file); os.IsNotExist(err) {
		fmt.Println("Creating new settings file at", s.file)
		s.Config = Config{}
		s.Config.Repositories = []Repository{}
		s.Config.Backups = []Backup{}
		s.Config.Mounts = []Mount{}
		s.Config.Schedules = []Schedule{}
		s.Save(s.Config)
	} else {
		fmt.Println("Loading settings from existing file", s.file)
		s.Config = s.readFile()
	}

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
	if file, err := os.Open(s.file); err == nil {
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

func (s *Settings) Save(data Config) error {
	s.Config = data
	fmt.Println("Saving settings")
	if str, err := json.MarshalIndent(s.Config, " ", " "); err == nil {
		fmt.Println("Settings saved")
		if err := os.WriteFile(s.file, str, 0644); err != nil {
			fmt.Println("error", err)
			return err
		}
	} else {
		fmt.Println("error", err)
		return err
	}
	return nil
}
