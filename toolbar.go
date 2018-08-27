package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
)

type toolbar struct {
	*gui.Panel
}

func (app *tdApp) buildToolbar() (tb *toolbar) {
	tb = new(toolbar)
	tb.Panel = gui.NewPanel(videoWidth, 28)
	tb.SetContentWidth(videoWidth)
	//tb.SetBorders(1, 1, 1, 1)
	tb.SetMargins(1, 1, 1, 1)

	hbl := gui.NewHBoxLayout()
	hbl.SetSpacing(4)
	tb.SetLayout(hbl)

	stopBtn := gui.NewButton("")
	stopBtn.SetIcon(icon.Pause)
	stopBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		drone.Hover()
	})
	tb.Add(stopBtn)

	cameraBtn := gui.NewButton("")
	cameraBtn.SetIcon(icon.CameraAlt)
	cameraBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		drone.TakePicture()
	})
	tb.Add(cameraBtn)

	setHomeBtn := gui.NewButton("")
	setHomeBtn.SetIcon(icon.AddLocation)
	setHomeBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		drone.SetHome()
	})
	tb.Add(setHomeBtn)

	goHomeBtn := gui.NewButton("")
	goHomeBtn.SetIcon(icon.Place)
	goHomeBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		if drone.IsHomeSet() {
			drone.AutoFlyToXY(0, 0)
		} else {
			alertDialog(app.mainPanel, warningSev, "Home position not set")
		}
	})
	tb.Add(goHomeBtn)

	return tb
}
