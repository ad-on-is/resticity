package internal

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
		os.Mkdir(filepath.Dir(s.file), 0755)
		os.Create(s.file)
		s.Init()
	} else {
		log.Info("Loading existing settings", "file", s.file)
		if s.FileEmpty() {
			s.Init()
		} else {
			s.Config = s.readFile()
		}
	}

	s.mux = sync.Mutex{}

	return s
}

func (s *Settings) Init() {
	log.Info("Initializing new settings", "file", s.file)
	s.Config = s.freshConfig()
	s.Save(s.Config)
}

func (s *Settings) freshConfig() Config {
	c := Config{}
	c.Repositories = []Repository{}
	c.Backups = []Backup{}
	c.Mounts = []Mount{}
	c.Schedules = []Schedule{}
	c.AppSettings = AppSettings{
		Theme: "auto",
		Notifications: AppSettingsNotifications{
			OnScheduleError:   true,
			OnScheduleSuccess: true,
			OnScheduleStart:   true,
		},
		Hooks: AppSettingsHooks{
			OnScheduleError:   "",
			OnScheduleSuccess: "",
			OnScheduleStart:   "",
		},
		PreserveErrorLogsDays: 7,
	}
	return c
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

func (c *Config) GetRepositoryById(id string) *Repository {
	for _, r := range c.Repositories {
		if r.Id == id {
			return &r
		}
	}
	return nil
}

func (c *Config) GetScheduleObject(s *Schedule) ScheduleObject {
	so := ScheduleObject{}
	so.Schedule = *s
	so.Backup = c.GetBackupById(s.BackupId)
	so.FromRepository = c.GetRepositoryById(s.FromRepositoryId)
	so.ToRepository = c.GetRepositoryById(s.ToRepositoryId)
	return so
}

func (c *Config) GetBackupById(id string) *Backup {
	for _, b := range c.Backups {
		if b.Id == id {
			return &b
		}
	}
	return nil
}

func (s *Settings) FileEmpty() bool {
	data, err := os.ReadFile(s.file)
	if err != nil {
		log.Error("file empty", "err", err)
		return true
	}
	return len(data) == 0
}

func (s *Settings) readFile() Config {
	s.mux.Lock()
	defer s.mux.Unlock()
	data := s.freshConfig()
	if file, err := os.Open(s.file); err == nil {
		if str, err := io.ReadAll(file); err == nil {
			if err := json.Unmarshal([]byte(str), &data); err != nil {
				log.Error("settings: unmarshal", "err", err)
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
