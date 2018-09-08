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
		alert := gtk.NewMessageDialog(win, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE,
			err.Error())
		alert.SetTitle(appName)
		alert.Run()
		alert.Destroy()
	}
	log.Printf("Saved %d photos", n)
}
