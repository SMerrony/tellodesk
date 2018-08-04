package main

import (
	"fmt"
	"strings"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/util/application"
)

type tdApp struct {
	*application.Application
	settingsLoaded                                                   bool
	settings                                                         settingsT
	menuBar                                                          *gui.Menu
	mainPanel                                                        *gui.Panel
	fileMenu, droneMenu, flightMenu, videoMenu, imagesMenu, helpMenu *gui.Menu
	connectItem, disconnectItem                                      *gui.MenuItem
	panel                                                            *gui.Panel
	label                                                            *gui.Label
}

func (app *tdApp) setup() {
	app.Gui().SetLayout(gui.NewVBoxLayout())
	// most stuff happens on the main panel
	app.mainPanel = gui.NewPanel(prefWidth, prefHeight)
	app.Gui().Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		app.mainPanel.SetWidth(app.Gui().ContentWidth())
		app.mainPanel.SetHeight(app.Gui().ContentHeight())
	})
	app.Gui().Add(app.mainPanel)

	// load any saved settings now as they may affect the gui
	var err error
	app.settings, err = loadSettings(appSettingsFile)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find") {
			alertDialog(app, warningSev, "Could not open settings file\n\n"+appSettingsFile+"\n\n"+
				"This is normal on a first run,\nor until you have saved your settings")
		} else {
			alertDialog(app, warningSev, err.Error())
		}
		app.settingsLoaded = false
		app.Log().Info("Error loading saved settings: %v", err)
	} else {
		app.settingsLoaded = true
	}

	app.buildMenu()
	app.mainPanel.Add(app.menuBar)
	app.Gui().SetName(appName)
	app.Subscribe(application.OnQuit, app.exitNicely) // catch main window being closed
}

func (app *tdApp) buildMenu() {
	app.menuBar = gui.NewMenuBar()
	app.fileMenu = gui.NewMenu()
	app.fileMenu.AddOption("Settings").Subscribe(gui.OnClick, app.settingsDialog)
	app.fileMenu.AddSeparator()
	//app.fileMenu.AddOption("Exit").SetId("exit").Subscribe(gui.OnClick, func(s string, i interface{}) { app.Quit() })
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
	app.helpMenu.AddOption("About").Subscribe(gui.OnClick, app.aboutCB)
	app.menuBar.AddMenu("Help", app.helpMenu)
}

func (app *tdApp) exitNicely(s string, i interface{}) {
	app.UnsubscribeID(application.OnQuit, nil) // prevent this being called again due to window app.Quit subscription
	app.Log().Info("Tidying-up and exiting")
	app.Quit()
}

func (app *tdApp) aboutCB(s string, i interface{}) {
	alertDialog(
		app,
		infoSev,
		fmt.Sprintf("Version: %s\n\nAuthor: %s\n\nCopyright: %s\n\nDisclaimer: %s", appVersion, appAuthor, appCopyright, appDisclaimer))
}
