package main

import (
	"fmt"
	"time"
)

func (app *tdApp) takePhotoCB(s string, i interface{}) {
	drone.TakePicture()
}

func (app *tdApp) saveAllPhotosCB(s string, i interface{}) {
	// TODO - use prefered directory from settings
	n, err := drone.SaveAllPics(fmt.Sprintf("tello_pic_%s", time.Now().Format("2006Jan2150405")))
	if err != nil {
		app.Log().Info("Error saving photos: %s", err.Error())
		alertDialog(app.mainPanel, errorSev, err.Error())
	}
	app.Log().Info("Saved %d photos", n)
}
