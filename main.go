package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/webfraggle/mbd-cli/internal/api"
	"github.com/webfraggle/mbd-cli/internal/config"
	"github.com/webfraggle/mbd-cli/internal/spawn"
	"github.com/webfraggle/mbd-cli/internal/ui"
)

func main() {
	// No args → open config UI
	if len(os.Args) == 1 {
		ui.Run()
		return
	}

	bg := flag.Bool("bg", false, "")
	next := flag.Bool("next", false, "")
	prev := flag.Bool("prev", false, "")
	setTime := flag.String("setTime", "", "")
	setTrain1 := flag.String("setTrain1", "", "")
	setTrain2 := flag.String("setTrain2", "", "")
	setTrain3 := flag.String("setTrain3", "", "")
	image := flag.String("image", "", "")
	gleis := flag.String("gleis", "", "")
	conf := flag.String("conf", "", "")
	timeout := flag.Int("timeout", 30000, "")

	flag.CommandLine.Usage = func() {} // suppress default usage output
	flag.Parse()

	// Foreground: spawn detached background copy and exit immediately
	if !*bg {
		spawn.Detach()
		return
	}

	// Background worker: set up logging, then execute the HTTP call
	setupLog()
	log.Println("Args:", os.Args)

	cfg, err := config.Load(*conf)
	if err != nil {
		log.Println("ERROR loading config:", err)
		os.Exit(1)
	}

	path := "GleisA"
	if len(*gleis) > 0 && strings.ToLower(string((*gleis)[len(*gleis)-1])) == "b" {
		path = "GleisB"
	}

	log.Println("Endpoint:", cfg.Endpoint, "Path:", path)

	client := api.NewClient(cfg.Endpoint, *timeout)

	if *next {
		log.Println("skip next")
		if err := client.SkipNext(path); err != nil {
			log.Println("ERROR:", err)
			os.Exit(1)
		}
		log.Println("SUCCESS")
		return
	}

	if *prev {
		log.Println("skip prev")
		if err := client.SkipPrev(path); err != nil {
			log.Println("ERROR:", err)
			os.Exit(1)
		}
		log.Println("SUCCESS")
		return
	}

	if *setTime != "" {
		log.Println("setTime:", *setTime)
		if err := client.SetTime(path, *setTime); err != nil {
			log.Println("ERROR:", err)
			os.Exit(1)
		}
		log.Println("SUCCESS")
		return
	}

	if *image != "" {
		log.Println("showImage:", *image)
		if err := client.ShowImage(path, *image); err != nil {
			log.Println("ERROR:", err)
			os.Exit(1)
		}
		log.Println("SUCCESS")
		return
	}

	var trains []api.TrainSlot
	if *setTrain1 != "" {
		trains = append(trains, api.ParseTrain(1, *setTrain1, path))
	}
	if *setTrain2 != "" {
		trains = append(trains, api.ParseTrain(2, *setTrain2, path))
	}
	if *setTrain3 != "" {
		trains = append(trains, api.ParseTrain(3, *setTrain3, path))
	}

	if len(trains) > 0 {
		log.Printf("setTrains: %d train(s)", len(trains))
		if err := client.SetTrains(trains); err != nil {
			log.Println("ERROR:", err)
			os.Exit(1)
		}
		log.Println("SUCCESS")
	}
}

func setupLog() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	logPath := filepath.Join(filepath.Dir(exe), "debug.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	log.SetOutput(io.MultiWriter(f, os.Stderr))
	fmt.Fprintln(f) // newline separator between runs
}
