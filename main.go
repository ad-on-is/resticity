package main

import (
	"bytes"
	"embed"
	"flag"

	"github.com/ad-on-is/resticity/internal"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/.output/public
var assets embed.FS

func main() {
	log.SetLevel(log.DebugLevel)
	flagConfigFile := ""
	flagBackground := false
	flag.StringVar(&flagConfigFile, "config", "", "Specify a config file")
	flag.StringVar(&flagConfigFile, "c", "", "Specify a config file")
	flag.BoolVar(&flagBackground, "background", false, "Run in background mode")
	flag.BoolVar(&flagBackground, "b", false, "Run in background mode")
	flag.Parse()
	log.Info("settings file", flagConfigFile)
	errb := bytes.NewBuffer([]byte{})
	outb := bytes.NewBuffer([]byte{})
	outputChan := make(chan internal.ChanMsg)
	settings := internal.NewSettings(flagConfigFile)
	restic := internal.NewRestic(errb, outb, settings)
	if scheduler, err := internal.NewScheduler(settings, restic, &outputChan); err == nil {
		(scheduler).RescheduleBackups()

		go internal.RunServer(scheduler, restic, settings, errb, outb, &outputChan)
		Desktop(scheduler, restic, settings, flagBackground)
	} else {
		log.Error("Init scheduler", "error", err)
	}

}

func Desktop(
	scheduler *internal.Scheduler,
	restic *internal.Restic,
	settings *internal.Settings,
	isHidden bool,
) {
	// Create an instance of the app structure
	app := NewApp(restic, scheduler, settings, &assets)
	// Create application with options
	err := wails.Run(&options.App{
		Title:             "resticity",
		Width:             1024,
		Height:            768,
		HideWindowOnClose: true,
		StartHidden:       isHidden,

		LogLevel: logger.ERROR,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
