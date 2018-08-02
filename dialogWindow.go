package main

import (
	"github.com/g3n/engine/gui"
)

type severityType int

const (
	infoSev severityType = iota
	warningSev
	errorSev
	criticalSev
)

const (
	dlgWidth, dlgHeight = 300, 300
)

func alertDialog(sev severityType, msg string) *gui.Window {
	win := gui.NewWindow(dlgWidth, dlgHeight)
	win.SetResizable(false)
	win.SetPaddings(4, 4, 4, 4)

	var titleStr string
	switch sev {
	case infoSev:
		titleStr = "Information"
	case warningSev:
		titleStr = "Warning"
	case errorSev:
		titleStr = "Error"
	case criticalSev:
		titleStr = "Critical Error"
	}
	win.SetTitle(titleStr)

	lay := gui.NewVBoxLayout()
	lay.SetSpacing(4)
	win.SetLayout(lay)

	msgLab := gui.NewImageLabel(msg)
	msgLab.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 2, AlignH: gui.AlignWidth})
	win.Add(msgLab)

	ok := gui.NewButton("OK")
	ok.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignCenter})
	ok.Subscribe(gui.OnClick, func(e string, ev interface{}) { win.SetVisible(false) })
	win.Add(ok)

	return win
}
