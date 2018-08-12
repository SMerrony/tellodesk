package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
)

type toolbar struct {
	*gui.Panel
}

func buildToolbar(parent *gui.Panel) (tb *toolbar) {
	tb = new(toolbar)
	tb.Panel = gui.NewPanel(parent.Width(), 28)
	//tb.SetBorders(1, 1, 1, 1)
	tb.SetMargins(1, 1, 1, 1)

	hbl := gui.NewHBoxLayout()
	hbl.SetSpacing(4)
	tb.SetLayout(hbl)

	stopBtn := gui.NewButton("")
	stopBtn.SetIcon(icon.Pause)
	//stopBtn.SubscrBtne(gui.OnCursorEnter, tooltip)
	tb.Add(stopBtn)

	cameraBtn := gui.NewButton("")
	// cameraBtn.SetIcon(icon.Warning)
	cameraBtn.SetIcon(icon.CameraAlt)
	tb.Add(cameraBtn)

	setHomeBtn := gui.NewButton("")
	// setHomeBtn.SetIcon(icon.Stop)
	setHomeBtn.SetIcon(icon.AddLocation)
	tb.Add(setHomeBtn)

	goHome3Btn := gui.NewButton("")
	goHome3Btn.SetIcon(icon.Place)
	goHome3Btn.SetEnabled(false)
	tb.Add(goHome3Btn)

	next4Btn := gui.NewButton("")
	next4Btn.SetIcon(icon.Error)
	tb.Add(next4Btn)

	return tb
}
