package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/ad-on-is/resticity/internal"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/.output/public
var assets embed.FS

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

	r.Scheduler.Assets = &assets
	if err == nil {
		(r.Scheduler).RescheduleBackups()
		go internal.RunServer(
			r.Scheduler,
			r.Restic,
			r.Settings,
			&r.OutputChan,
			&r.ErrorChan,
			Version,
			Build,
		)
		Desktop(r.Scheduler, r.Restic, r.Settings, r.FlagArgs.Background)
	} else {
		log.Error("Resticity failed to start", "error", err)
		os.Exit(1)
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
