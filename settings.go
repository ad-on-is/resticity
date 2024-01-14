package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/adrg/xdg"
)

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
