package main

import (
	"fmt"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/util/application"
)

type tdApp struct {
	*application.Application
	menuBar                                                          *gui.Menu
	fileMenu, droneMenu, flightMenu, videoMenu, imagesMenu, helpMenu *gui.Menu
	connectItem, disconnectItem                                      *gui.MenuItem
	panel                                                            *gui.Panel
	label                                                            *gui.Label
}

func (app *tdApp) setup() {
	app.Gui().SetLayout(gui.NewVBoxLayout())

	app.buildMenu()
	app.Gui().Add(app.menuBar)

}

func (app *tdApp) buildMenu() {
	app.menuBar = gui.NewMenuBar()
	app.menuBar.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignWidth})

	app.fileMenu = gui.NewMenu()
	app.fileMenu.AddOption("Exit").SetId("exit").Subscribe(gui.OnClick, app.exitNicely)
	app.menuBar.AddMenu("File", app.fileMenu)

	//app.menuBar.AddSeparator()

	app.droneMenu = gui.NewMenu()
	app.connectItem = app.droneMenu.AddOption("Connect")
	app.connectItem.Subscribe(gui.OnClick, app.connectCB)
	app.disconnectItem = app.droneMenu.AddOption("Disconnect")
	app.disconnectItem.SetEnabled(false).Subscribe(gui.OnClick, app.diconnectCB)
	app.menuBar.AddMenu("Drone", app.droneMenu)

	app.flightMenu = gui.NewMenu()
	app.flightMenu.AddOption("Take-off").Subscribe(gui.OnClick, app.takeoffCB)
	app.flightMenu.AddOption("Throw Take-off").Subscribe(gui.OnClick, app.throwTakeoffCB)
	app.flightMenu.AddOption("Land").Subscribe(gui.OnClick, app.landCB)
	app.flightMenu.AddOption("Palm Land").Subscribe(gui.OnClick, app.palmLandCB)
	app.flightMenu.AddSeparator()
	app.flightMenu.AddOption("Sports (Fast) Mode")
	app.menuBar.AddMenu("Flight", app.flightMenu)

	app.videoMenu = gui.NewMenu()
	app.videoMenu.AddOption("Start Video View")
	app.videoMenu.AddOption("Stop Video View")
	app.videoMenu.AddSeparator()
	app.videoMenu.AddOption("Record Video")
	app.menuBar.AddMenu("Video", app.videoMenu)

	app.imagesMenu = gui.NewMenu()
	app.imagesMenu.AddOption("Take Photo")
	app.imagesMenu.AddOption("Save Photo(s)").SetEnabled(false)
	app.menuBar.AddMenu("Images", app.imagesMenu)

	app.helpMenu = gui.NewMenu()
	app.helpMenu.AddOption("Online Help")
	app.helpMenu.AddSeparator()
	app.helpMenu.AddOption("About")
	app.menuBar.AddMenu("Help", app.helpMenu)
}

func (app *tdApp) exitNicely(s string, i interface{}) {
	fmt.Println("Tidying-up and exiting")
	app.Quit()
}
