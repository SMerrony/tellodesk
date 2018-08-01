package main

import "github.com/g3n/engine/gui"

type severityType int

const (
	infoSev severityType = iota
	warningSev
	errorSev
	criticalSev
)

const (
	alDwidth, alDheight = 300, 300
)

type alertDialog struct {
	gui.Panel
	severity   severityType
	title, msg *gui.ImageLabel
	ok         *gui.Button
}

func runAlert(sev severityType, msg string) *alertDialog {

	dlg := new(alertDialog)
	dlg.Initialize(alDwidth, alDheight)
	dlg.SetBorders(2, 2, 2, 2)
	dlg.SetPaddings(4, 4, 4, 4)

	lay := gui.NewVBoxLayout()
	lay.SetSpacing(4)
	dlg.SetLayout(lay)

	var typeStr string
	switch sev {
	case infoSev:
		typeStr = "Information"
	case warningSev:
		typeStr = "Warning"
	case errorSev:
		typeStr = "Error"
	case criticalSev:
		typeStr = "Critical Error"
	}
	dlg.title = gui.NewImageLabel(typeStr)
	dlg.title.SetBorders(1, 1, 1, 1)
	dlg.title.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignWidth})
	dlg.Add(dlg.title)

	dlg.msg = gui.NewImageLabel(msg)
	dlg.msg.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 2, AlignH: gui.AlignWidth})
	dlg.Add(dlg.msg)

	dlg.ok = gui.NewButton("OK")
	dlg.ok.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignCenter})
	dlg.ok.Subscribe(gui.OnClick, func(e string, ev interface{}) { dlg.SetVisible(false) })
	dlg.Add(dlg.ok)

	return dlg
}
