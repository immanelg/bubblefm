package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const version = "v0.1.0"

func getInitialCwd() string {
	var err error
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not get PWD")
		cwd, err = os.UserHomeDir()
		if err != nil {
			panic("could not get current directory and user home directory")
		}

	}
	return cwd
}

func main() {
	configFile := flag.String("config", "", "path to config file")
	versionFlag := flag.Bool("version", false, "print version and exit")
	logpath := flag.String("log", "", "print log to file")

	flag.Usage = func() {
		fmt.Printf("bubblefm %v, a simple file manager\n", version)
		fmt.Println("usage: bubblefm [OPTIONS...] [PATH]")
		fmt.Println("options:")
		fmt.Println("\t-help: print help and exit")
		fmt.Println("\t-version: print version and exit")
		fmt.Println("\t-log: logging file")
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	var cwd string
	if flag.NArg() == 1 {
		cwd = flag.Arg(0)
	} else {
		cwd = getInitialCwd()
	}

	if *logpath != "" {
		file, err := os.OpenFile(*logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		log.SetOutput(file)
	} else {
		log.SetOutput(io.Discard)
	}

	if *configFile == "" {
		configdir, err := os.UserConfigDir()
		if err != nil {
			panic("could not get user config dir")
		}
		*configFile = filepath.Join(configdir, "bubblefm", "config")
	}
	config := initConfig()
	_ = config.source(*configFile) // who cares about the errors? it's a user's problem.

	lipgloss.SetColorProfile(termenv.ANSI256)
	m := newModel(cwd, config)
	program := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}
