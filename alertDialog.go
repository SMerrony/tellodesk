package main

import (
	"log"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/math32"
)

type severityType int

const (
	infoSev severityType = iota
	warningSev
	errorSev
	criticalSev
)

const (
	dlgWidth, dlgHeight float32 = 300.0, 200.0
)

type alert struct {
	*gui.Window
}

func alertDialog(parent *gui.Panel, sev severityType, msg string) (win *alert) {
	win = new(alert)
	win.Window = gui.NewWindow(dlgWidth, dlgHeight)
	win.SetResizable(false)
	win.SetPaddings(4, 4, 4, 4)
	win.SetColor(math32.NewColor("Gray"))

	var iconStr string
	titleStr := appName + " - "
	switch sev {
	case infoSev:
		titleStr += "Information"
		iconStr = string(icon.Info)
	case warningSev:
		titleStr += "Warning"
		iconStr = string(icon.Warning)
	case errorSev:
		titleStr += "Error"
		iconStr = string(icon.Error)
	case criticalSev:
		titleStr += "Critical Error"
		iconStr = string(icon.Error)
	}
	win.SetTitle(titleStr)

	lay := gui.NewVBoxLayout()
	lay.SetSpacing(4)
	win.SetLayout(lay)

	msgLab := gui.NewImageLabel(msg)
	msgLab.SetIcon(iconStr)
	msgLab.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 2, AlignH: gui.AlignWidth})
	win.Add(msgLab)

	ok := gui.NewButton("OK")
	ok.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignCenter})
	ok.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		parent.Root().SetModal(nil)
		parent.Remove(win)
	})
	win.Add(ok)

	win.SetCloseButton(false)

	win.Subscribe(gui.OnClick, func(e string, ev interface{}) { log.Println("What?") })

	parent.Add(win)
	win.SetPosition(parent.Width()/2-dlgWidth/2, parent.Height()/2-dlgHeight/2)
	parent.Root().SetModal(win)

	return win
}
