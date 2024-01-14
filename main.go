package main

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/.output/public
var assets embed.FS

type ChanMsg struct {
	Id       string   `json:"id"`
	Out      string   `json:"out"`
	Schedule Schedule `json:"schedule"`
}

func main() {
	errb := bytes.NewBuffer([]byte{})
	outb := bytes.NewBuffer([]byte{})
	outputChan := make(chan ChanMsg)
	settings := NewSettings()
	restic := NewRestic(errb, outb, settings)
	if scheduler, err := NewScheduler(settings, restic, &outputChan); err == nil {
		(scheduler).RescheduleBackups()
		go RunServer(scheduler, restic, settings, errb, outb, &outputChan)
		Desktop(scheduler, restic, settings)
	} else {
		fmt.Println("SCHEDULER ERROR", err)
	}

}

func Desktop(scheduler *Scheduler, restic *Restic, settings *Settings) {
	// Create an instance of the app structure
	app := NewApp(restic, scheduler, settings)
	// Create application with options
	err := wails.Run(&options.App{
		Title:             "resticity",
		Width:             1024,
		Height:            768,
		HideWindowOnClose: true,
		LogLevel:          logger.ERROR,
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
