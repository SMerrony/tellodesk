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
	prefWidth, prefHeight = 1300, 800
)

var (
	drone      tello.Tello
	stickChan  chan<- tello.StickMessage
	jsStopChan chan bool
)

func main() {
	var err error
	td := new(tdApp)
	td.Application, err = application.Create(application.Options{
		Title:       appName,
		Width:       prefWidth,
		Height:      prefHeight,
		TargetFPS:   30,
		EnableFlags: true,
	})
	if err != nil {
		log.Fatalf("Error creating application: %v", err)
	}
	jsStopChan = make(chan bool) // not buffered
	td.setup()

	td.Run()
}
