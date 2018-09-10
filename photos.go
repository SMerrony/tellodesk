package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/mattn/go-gtk/gtk"
)

func takePhotoCB() {
	drone.TakePicture()
}

func saveAllPhotosCB() {
	n, err := drone.SaveAllPics(fmt.Sprintf("%s%ctello_pic_%s",
		settings.DataDir, filepath.Separator, time.Now().Format("2006Jan2150405")))
	if err != nil {
		log.Printf("Error saving photos: %s", err.Error())
		messageDialog(win, gtk.MESSAGE_ERROR, err.Error())
	}
	log.Printf("Saved %d photos", n)
}
