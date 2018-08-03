package main

import (
	"github.com/g3n/engine/gui"
)

const (
	dialogTitle                           = appName + " Settings"
	settingsWidth, settingsHeight float32 = 500, 200
)

func (app *tdApp) settingsDialog(s string, i interface{}) {
	win := gui.NewWindow(settingsWidth, settingsHeight)
	win.SetTitle(dialogTitle)
	win.SetCloseButton(false)
	//win.SetPaddings(8, 8, 8, 8)

	lay := gui.NewGridLayout(3)
	lay.SetAlignH(gui.AlignCenter)
	lay.SetExpandH(true)
	win.SetLayout(lay)

	win.Add(gui.NewLabel(""))
	win.Add(gui.NewLabel("Detected"))
	win.Add(gui.NewLabel("Type"))

	win.Add(gui.NewLabel("Joystick"))
	dDrop := gui.NewDropDown(200, gui.NewImageLabel(""))
	dDrop.SetMargins(3, 3, 3, 3)
	found := listJoysticks()
	for _, j := range found {
		dDrop.Add(gui.NewImageLabel(j.Name))
	}
	win.Add(dDrop)
	tDrop := gui.NewDropDown(150, gui.NewImageLabel(""))
	tDrop.SetMargins(3, 3, 3, 3)
	known := listKnownJoystickTypes()
	for _, k := range known {
		tDrop.Add(gui.NewImageLabel(k.Name))
	}
	win.Add(tDrop)

	// empty row...
	win.Add(gui.NewLabel(""))
	win.Add(gui.NewLabel(""))
	win.Add(gui.NewLabel(""))

	// buttons...
	win.Add(gui.NewLabel(""))
	cancel := gui.NewButton("Cancel")
	cancel.SetBorders(1, 1, 1, 1)
	cancel.SetPaddings(3, 3, 3, 3)
	cancel.SetMargins(3, 3, 3, 3)
	ok := gui.NewButton("OK")
	ok.SetBorders(1, 1, 1, 1)
	ok.SetPaddings(3, 3, 3, 3)
	ok.SetMargins(3, 3, 3, 3)
	win.Add(cancel)
	cancel.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		app.Log().Info("Cancelled")
		app.Gui().Root().SetModal(nil)
		app.Gui().Root().Remove(win)
	})
	win.Add(ok)

	root := app.Gui().Root()
	root.Add(win)
	win.SetPosition(root.Width()/2-settingsWidth/2, root.Height()/2-settingsHeight/2)
	root.SetModal(win)
}
