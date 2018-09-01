package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
)

func (app *tdApp) buildMenu() {

	app.menuBar = gui.NewMenuBar()

	fileMenu := gui.NewMenu()
	settings := fileMenu.AddOption("Settings")
	settings.SetIcon(icon.Settings)
	settings.Subscribe(gui.OnClick, app.settingsCB)

	fileMenu.AddSeparator()

	ex := fileMenu.AddOption("Exit")
	ex.SetId("exit")
	ex.SetIcon(icon.Close)
	ex.Subscribe(gui.OnClick, app.exitNicely)

	app.menuBar.AddMenu("File ", fileMenu)

	droneMenu := gui.NewMenu()
	app.connectItem = droneMenu.AddOption("Connect")
	app.connectItem.SetIcon(icon.Sync)
	app.connectItem.Subscribe(gui.OnClick, app.connectCB)
	app.disconnectItem = droneMenu.AddOption("Disconnect")
	app.disconnectItem.SetIcon(icon.SyncDisabled)
	app.disconnectItem.SetEnabled(false).Subscribe(gui.OnClick, app.disconnectCB)

	app.flightSubMenu = gui.NewMenu()

	to := app.flightSubMenu.AddOption("Take-off")
	to.SetIcon(icon.FlightTakeoff)
	to.Subscribe(gui.OnClick, app.takeoffCB)
	tto := app.flightSubMenu.AddOption("Throw Take-off")
	tto.SetIcon(icon.ThumbUp)
	tto.Subscribe(gui.OnClick, app.throwTakeoffCB)
	lnd := app.flightSubMenu.AddOption("Land")
	lnd.SetIcon(icon.FlightLand)
	lnd.Subscribe(gui.OnClick, app.landCB)
	plnd := app.flightSubMenu.AddOption("Palm Land")
	plnd.SetIcon(icon.PanTool)
	plnd.Subscribe(gui.OnClick, app.palmLandCB)
	app.flightSubMenu.AddSeparator()
	sm := app.flightSubMenu.AddOption("Sports (Fast) Mode")
	sm.SetIcon(icon.DirectionsRun)
	sm.Subscribe(gui.OnClick, app.nyi)

	droneMenu.AddMenu("Flight", app.flightSubMenu)
	app.flightSubMenu.SetEnabled(false)

	app.menuBar.AddMenu(" Drone ", droneMenu)

	app.trackMenu = gui.NewMenu()
	ct := app.trackMenu.AddOption("Clear Track")
	ct.SetIcon(icon.Delete)
	ct.Subscribe(gui.OnClick, app.nyi)
	et := app.trackMenu.AddOption("Export Current Track as CSV")
	et.SetIcon(icon.Save)
	et.Subscribe(gui.OnClick, app.exportTrackCB)
	app.importTrackItem = app.trackMenu.AddOption("Import CSV Track")
	app.importTrackItem.SetIcon(icon.Input)
	app.importTrackItem.Subscribe(gui.OnClick, app.importTrackCB)
	st := app.trackMenu.AddOption("Save Track as PNG")
	st.SetIcon(icon.Image)
	st.Subscribe(gui.OnClick, app.exportTrackImageCB)

	trackSubMenu := gui.NewMenu()

	app.tsmShowDrone = trackSubMenu.AddOption("Show Drone Positions")
	app.tsmShowDrone.SetIcon(icon.CheckBox)
	app.trackShowDrone = true
	app.tsmShowDrone.Subscribe(gui.OnClick, app.trackShowDroneCB)
	app.tsmShowPath = trackSubMenu.AddOption("Show Track Path")
	app.tsmShowPath.SetIcon(icon.CheckBox)
	app.trackShowPath = true
	app.tsmShowPath.Subscribe(gui.OnClick, app.trackShowPathCB)

	app.trackMenu.AddMenu("Display", trackSubMenu)

	app.menuBar.AddMenu(" Track ", app.trackMenu)

	videoMenu := gui.NewMenu()
	app.recordVideoItem = videoMenu.AddOption("Record Video")
	app.recordVideoItem.SetIcon(icon.Videocam)
	app.recordVideoItem.Subscribe(gui.OnClick, app.recordVideoCB)
	app.stopRecordingItem = videoMenu.AddOption("Stop Recording")
	app.stopRecordingItem.SetIcon(icon.VideocamOff)
	app.stopRecordingItem.Subscribe(gui.OnClick, app.stopRecordingCB)
	app.stopRecordingItem.SetEnabled(false)
	app.menuBar.AddMenu(" Video ", videoMenu)

	app.imagesMenu = gui.NewMenu()
	tp := app.imagesMenu.AddOption("Take Photo")
	tp.SetIcon(icon.CameraAlt)
	tp.Subscribe(gui.OnClick, app.takePhotoCB)
	sp := app.imagesMenu.AddOption("Save Photo(s)")
	sp.SetIcon(icon.Save)
	sp.Subscribe(gui.OnClick, app.saveAllPhotosCB)
	app.menuBar.AddMenu(" Images ", app.imagesMenu)

	helpMenu := gui.NewMenu()
	oh := helpMenu.AddOption("Online Help")
	oh.SetIcon(icon.Help)
	oh.Subscribe(gui.OnClick, app.onlineHelpCB)
	helpMenu.AddSeparator()
	ab := helpMenu.AddOption("About")
	ab.SetIcon(icon.Info)
	ab.Subscribe(gui.OnClick, app.aboutCB)
	app.menuBar.AddMenu(" Help", helpMenu)

	app.menuBar.SetWidth(videoWidth)
}

func (app *tdApp) enableFlightMenus() {
	app.disconnectItem.SetEnabled(true)
	app.connectItem.SetEnabled(false)
	app.flightSubMenu.SetEnabled(true)
	app.importTrackItem.SetEnabled(false)
}

func (app *tdApp) disableFlightMenus() {
	app.disconnectItem.SetEnabled(false)
	app.connectItem.SetEnabled(true)
	app.flightSubMenu.SetEnabled(false)
	app.importTrackItem.SetEnabled(true)
}
