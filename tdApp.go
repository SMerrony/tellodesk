package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/util/application"
)

type tdApp struct {
	*application.Application
	// gui                tdGUI
	menuBar            *gui.Menu
	fileMenu, helpMenu *gui.Menu
	panel              *gui.Panel
	label              *gui.Label
}

func (app *tdApp) setup() {
	app.Gui().SetLayout(gui.NewVBoxLayout())
	app.menuBar = gui.NewMenuBar()
	app.menuBar.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignWidth})
	app.Gui().Add(app.menuBar)
	app.fileMenu = gui.NewMenu()
	app.fileMenu.AddOption("Exit").SetId("exit").Subscribe(gui.OnClick, func(string, interface{}) {
		app.Quit()
	})
	app.menuBar.AddMenu("File", app.fileMenu)

	//app.menuBar.AddSeparator()

	app.helpMenu = gui.NewMenu()
	app.helpMenu.AddSeparator()
	app.helpMenu.AddOption("About")
	app.menuBar.AddMenu("Help", app.helpMenu)
	//app.Gui().Add(app.gui.window)

}
