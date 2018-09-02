package main

import (
	"fmt"
	"path/filepath"
	"time"
)

func (app *tdApp) takePhotoCB(s string, i interface{}) {
	drone.TakePicture()
}

func (app *tdApp) saveAllPhotosCB(s string, i interface{}) {
	n, err := drone.SaveAllPics(fmt.Sprintf("%s%ctello_pic_%s",
		app.settings.DataDir, filepath.Separator, time.Now().Format("2006Jan2150405")))
	if err != nil {
		app.Log().Info("Error saving photos: %s", err.Error())
		alertDialog(app.mainPanel, errorSev, err.Error())
	}
	app.Log().Info("Saved %d photos", n)
}
