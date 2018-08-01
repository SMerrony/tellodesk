package main

import (
	"log"

	"github.com/g3n/engine/util/application"
)

const (
	appName               = "Tello Desktop"
	prefWidth, prefHeight = 800, 600
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
