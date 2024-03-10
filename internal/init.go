package internal

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func GetVersion() string {
	return "0.0.1"
}

type FlagArgs struct {
	ConfigFile string
	Help       bool
	Version    bool
	Background bool
}

type Resticity struct {
	FlagArgs   FlagArgs
	ErrB       *bytes.Buffer
	OutB       *bytes.Buffer
	OutputChan chan ChanMsg
	Settings   *Settings
	Restic     *Restic
	Scheduler  *Scheduler
}

func NewResticity() (Resticity, error) {
	flagArgs := ParseFlags()
	flagArgs.PrintVersionOrHelp()
	errb := bytes.NewBuffer([]byte{})
	outb := bytes.NewBuffer([]byte{})
	outputChan := make(chan ChanMsg)
	settings := NewSettings(flagArgs.ConfigFile)
	restic := NewRestic(errb, outb, settings)
	scheduler, err := NewScheduler(settings, restic, &outputChan)
	return Resticity{flagArgs, errb, outb, outputChan, settings, restic, scheduler}, err
}

func ParseFlags() FlagArgs {
	flagArgs := FlagArgs{ConfigFile: "", Help: false, Version: false}

	flag.StringVar(&flagArgs.ConfigFile, "config", "", "Specify a config file")
	flag.StringVar(&flagArgs.ConfigFile, "c", "", "Specify a config file")
	flag.BoolVar(&flagArgs.Background, "background", false, "Run in background mode")
	flag.BoolVar(&flagArgs.Background, "b", false, "Run in background mode")
	flag.BoolVar(&flagArgs.Help, "help", false, "Show help")
	flag.BoolVar(&flagArgs.Help, "h", false, "Show help")
	flag.BoolVar(&flagArgs.Version, "version", false, "Show version")
	flag.BoolVar(&flagArgs.Version, "v", false, "Show version")
	flag.Parse()

	return flagArgs
}

func (flagArgs *FlagArgs) PrintVersionOrHelp() {
	if flagArgs.Version {
		fmt.Println("resticity " + GetVersion())
		os.Exit(0)
	}

	if flagArgs.Help {
		fmt.Println("resticity " + GetVersion())
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func SetLogLevel() {
	l, err := log.ParseLevel(os.Getenv("RESTICITY_LOG_LEVEL"))
	if err == nil {
		log.SetLevel(l)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
