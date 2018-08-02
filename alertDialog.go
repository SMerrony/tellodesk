package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
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

func alertDialog(root *gui.Root, sev severityType, msg string) {
	win := gui.NewWindow(dlgWidth, dlgHeight)
	win.SetResizable(false)
	win.SetPaddings(4, 4, 4, 4)

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
		root.SetModal(nil)
		root.Remove(win)
	})
	win.Add(ok)

	win.SetCloseButton(false)

	root.Add(win)
	win.SetPosition(root.Width()/2-dlgWidth/2, root.Height()/2-dlgHeight/2)
	root.SetModal(win)
}
