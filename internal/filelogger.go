package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
)

func getPath() string {
	return filepath.Join(xdg.CacheHome, "resticity")
}

func getFile(t string) string {
	d := time.Now().Format("2006-01-02")
	f := filepath.Join(getPath(), t+"_"+d+".log")
	return f
}

func appendToFile(f string, m ChanMsg) {

	if _, err := os.Stat(f); os.IsNotExist(err) {
		_, err := os.Create(f)
		if err != nil {
			log.Error("filelogger: create "+f, "error", err)
			return
		}
	}

	d, err := json.Marshal(m)
	if err != nil {
		log.Error("filelogger: marshal "+f, "error", err)
		return
	}

	if err := WriteFile(f, []byte(d)); err != nil {
		log.Error("filelogger: write "+f, "error", err)
	}

}

func WriteFile(name string, data []byte) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(string(data) + "\n"); err != nil {
		return err
	}
	return nil
}

func NewFileLogger(outputChan *chan ChanMsg, errorChan *chan ChanMsg) {
	if _, err := os.Stat(getPath()); os.IsNotExist(err) {
		os.Mkdir(getPath(), 0755)
	}
	log.Info("filelogger", "path", getPath())
	for {
		select {
		case o := <-*outputChan:
			if os.Getenv("LOG_TO_FILE") == "true" {
				f := getFile("logs")
				appendToFile(f, o)
			}
			break
		case e := <-*errorChan:
			f := getFile("error")
			appendToFile(f, e)
			break
		}

	}
}
