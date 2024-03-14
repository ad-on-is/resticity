package internal

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

func RunHook(cmd string, obj ScheduleObject) {
	stat, err := os.Stat(cmd)
	if os.IsNotExist(err) {
		log.Error("Hook file not found", "file", cmd)
		return
	}
	if stat.Mode()&0111 == 0 {
		log.Error("Hook file not executable", "file", cmd)
		return
	}
	u, err := json.Marshal(obj)
	if err != nil {
		log.Error("Failed to marshal object", "error", err)
		return
	}
	s := string(u)
	execCmd := exec.Cmd{Path: cmd, Args: []string{cmd, s}}
	log.Debug(s)
	if err := execCmd.Start(); err != nil {
		log.Error("Hook failed", "file", cmd, "error", err)
		return
	}

	go execCmd.Wait()

}
