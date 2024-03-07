package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
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
		log.Info("Creating new settings file", "file", s.file)
		s.Config = Config{}
		s.Config.Repositories = []Repository{}
		s.Config.Backups = []Backup{}
		s.Config.Mounts = []Mount{}
		s.Config.Schedules = []Schedule{}
		s.Save(s.Config)
	} else {
		log.Info("Loading existing settings", "file", s.file)
		s.Config = s.readFile()
	}

	s.mux = sync.Mutex{}

	return s
}

func (s *Settings) SetLastRun(id string, error string) {
	for i, j := range s.Config.Schedules {
		if j.Id == id {
			log.Debug("save last run", "i", i, "id", id)
			s.Config.Schedules[i].LastRun = time.Now().Format(time.RFC3339)
			s.Config.Schedules[i].LastError = error

			s.Save(s.Config)
			break
		}
	}
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
	s.mux.Lock()
	defer s.mux.Unlock()
	data := Config{}
	if file, err := os.Open(s.file); err == nil {
		if str, err := io.ReadAll(file); err == nil {
			if err := json.Unmarshal([]byte(str), &data); err == nil {
				return data
			}
		}
	} else {
		log.Error("settings: read file", "err", err)
	}
	return data
}

func (s *Settings) Refresh() {
	s.Config = s.readFile()
}

func (s *Settings) Save(data Config) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Config = data
	log.Debug("Saving settings")
	if str, err := json.MarshalIndent(s.Config, " ", " "); err == nil {
		log.Info("Settings saved")
		if err := os.WriteFile(s.file, str, 0644); err != nil {
			log.Error("settings: write", "err", err)
			return err
		}
	} else {
		log.Error("settings: marshal indent", "err", err)
		return err
	}
	return nil
}
