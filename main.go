package main

import (
	"embed"
	"os"
	"strings"

	"github.com/ad-on-is/resticity/internal"
	"github.com/thoas/go-funk"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/.output/public
var assets embed.FS

func test(arr *[]string, n string) bool {
	return funk.Find(*arr, func(s string) bool {
		return strings.Contains(s, n)
	}) != nil
}

func main() {
	internal.SetLogLevel()
	r, err := internal.NewResticity()
	if err == nil {
		(r.Scheduler).RescheduleBackups()
		go internal.RunServer(
			r.Scheduler,
			r.Restic,
			r.Settings,
			&r.OutputChan,
			&r.ErrorChan,
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
