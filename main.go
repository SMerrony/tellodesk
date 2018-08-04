package main

import (
	"log"

	"github.com/SMerrony/tello"
	"github.com/g3n/engine/util/application"
)

const (
	appName               = "Tello Desktop"
	appVersion            = "0.1.0"
	appAuthor             = "S.Merrony"
	appCopyright          = "©2018 S.Merrony"
	appDisclaimer         = "The author(s) is/are in no way\nconnected with Ryze®."
	appSettingsFile       = "tellodesktop.yaml"
	prefWidth, prefHeight = 800, 600
)

var (
	drone     tello.Tello
	stickChan chan<- tello.StickMessage
)

func main() {
	td := new(tdApp)
	a, err := application.Create(application.Options{
		Title:       appName,
		Width:       prefWidth,
		Height:      prefHeight,
		TargetFPS:   30,
		EnableFlags: true,
	})
	if err != nil {
		log.Fatalf("Error creating application: %v", err)
	}

	td.Application = a

	td.setup()

	td.Run()
}
