package main

import "github.com/g3n/engine/gui"

type statusbar struct {
	*gui.Panel
	connectionLab *gui.Label
}

func buildStatusbar(parent *gui.Panel) (sb *statusbar) {
	sb = new(statusbar)
	sb.Panel = gui.NewPanel(parent.Width(), parent.Height())

	hbl := gui.NewHBoxLayout()
	hbl.SetSpacing(4)
	sb.SetLayout(hbl)

	sb.connectionLab = gui.NewLabel("Connection Status")
	sb.SetBorders(2, 2, 2, 2)
	sb.SetPaddings(3, 3, 3, 3)
	sb.Add(sb.connectionLab)

	return sb
}
