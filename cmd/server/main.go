package main

import (
	"os"

	"github.com/ad-on-is/resticity/internal"

	"github.com/charmbracelet/log"
)

func main() {
	internal.SetLogLevel()
	r, err := internal.NewResticity()
	if err == nil {
		(r.Scheduler).RescheduleBackups()
		internal.RunServer(r.Scheduler, r.Restic, r.Settings, r.ErrB, r.OutB, &r.OutputChan)

	} else {
		log.Error("Resticity failed to start", "error", err)
		os.Exit(1)
	}

}
