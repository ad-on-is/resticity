package main

import (
	"bytes"
	"flag"

	"github.com/ad-on-is/resticity/internal"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	flagConfigFile := ""
	flag.StringVar(&flagConfigFile, "config", "", "Specify a config file")
	flag.StringVar(&flagConfigFile, "c", "", "Specify a config file")
	flag.Parse()
	log.Info("settings file", flagConfigFile)
	errb := bytes.NewBuffer([]byte{})
	outb := bytes.NewBuffer([]byte{})
	outputChan := make(chan internal.ChanMsg)
	settings := internal.NewSettings(flagConfigFile)
	restic := internal.NewRestic(errb, outb, settings)
	if scheduler, err := internal.NewScheduler(settings, restic, &outputChan); err == nil {
		(scheduler).RescheduleBackups()
		internal.RunServer(scheduler, restic, settings, errb, outb, &outputChan)

	} else {
		log.Error("Init scheduler", "error", err)
	}

}
