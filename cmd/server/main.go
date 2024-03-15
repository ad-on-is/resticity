package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ad-on-is/resticity/internal"

	"github.com/charmbracelet/log"
)

var (
	Version string
	Build   string
)

func main() {
	internal.SetLogLevel()
	r, err := internal.NewResticity()
	if r.FlagArgs.Version {
		fmt.Println("resticity - version=" + Version + ", build=" + Build + "")
		os.Exit(0)
	}

	if r.FlagArgs.Help {
		fmt.Println("resticity " + Version + " (build " + Build + ")")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if err == nil {
		(r.Scheduler).RescheduleBackups()
		internal.RunServer(
			r.Scheduler,
			r.Restic,
			r.Settings,

			&r.OutputChan,
			&r.ErrorChan,
			Version,
			Build,
		)

	} else {
		log.Error("Resticity failed to start", "error", err)
		os.Exit(1)
	}

}
