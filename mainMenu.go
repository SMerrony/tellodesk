package main

import (
	"github.com/mattn/go-gtk/gtk"
)

type menuBarT struct {
	*gtk.MenuBar
	connectItem, disconnectItem   *gtk.MenuItem
	flightItem                    *gtk.MenuItem
	importTrackItem               *gtk.MenuItem
	recVidItem, stopRecVidItem    *gtk.MenuItem
	trackShowDrone, trackShowPath *gtk.CheckMenuItem
}

func buildMenu() (mb *menuBarT) {

	mb = new(menuBarT)
	mb.MenuBar = gtk.NewMenuBar()

	fileItem := gtk.NewMenuItemWithLabel("File")
	mb.Append(fileItem)
	fileMenu := gtk.NewMenu()
	fileItem.SetSubmenu(fileMenu)

	settings := gtk.NewMenuItemWithLabel("Settings")
	settings.Connect("activate", settingsCB)
	fileMenu.Append(settings)
	fileMenu.Append(gtk.NewSeparatorMenuItem())
	exitItem := gtk.NewMenuItemWithLabel("Exit")
	exitItem.Connect("activate", exitNicely)
	fileMenu.Append(exitItem)

	droneItem := gtk.NewMenuItemWithLabel("Drone")
	mb.Append(droneItem)
	droneMenu := gtk.NewMenu()
	droneItem.SetSubmenu(droneMenu)

	mb.connectItem = gtk.NewMenuItemWithLabel("Connect")
	mb.connectItem.Connect("activate", connectCB)
	droneMenu.Append(mb.connectItem)
	mb.disconnectItem = gtk.NewMenuItemWithLabel("Disconnect")
	mb.disconnectItem.Connect("activate", disconnectCB)
	droneMenu.Append(mb.disconnectItem)

	mb.flightItem = gtk.NewMenuItemWithLabel("Flight")
	droneMenu.Append(mb.flightItem)
	flightMenu := gtk.NewMenu()
	mb.flightItem.SetSubmenu(flightMenu)

	to := gtk.NewMenuItemWithLabel("Take-off")
	to.Connect("activate", takeoffCB)
	flightMenu.Append(to)
	tto := gtk.NewMenuItemWithLabel("Throw Take-off")
	tto.Connect("activate", throwTakeoffCB)
	flightMenu.Append(tto)
	li := gtk.NewMenuItemWithLabel("Land")
	li.Connect("activate", landCB)
	flightMenu.Append(li)
	pl := gtk.NewMenuItemWithLabel("Palm Land")
	pl.Connect("activate", palmLandCB)
	flightMenu.Append(pl)

	flightMenu.Append(gtk.NewSeparatorMenuItem())

	sm := gtk.NewMenuItemWithLabel("Sports (Fast) Mode")
	sm.Connect("activate", nyi)
	flightMenu.Append(sm)

	trackItem := gtk.NewMenuItemWithLabel("Track")
	mb.Append(trackItem)
	trackMenu := gtk.NewMenu()
	trackItem.SetSubmenu(trackMenu)

	// ct := gtk.NewMenuItemWithLabel("Clear Track")
	// ct.Connect("activate", nyi)
	// trackMenu.Append(ct)
	et := gtk.NewMenuItemWithLabel("Export Current Track as CSV")
	et.Connect("activate", exportTrackCB)
	trackMenu.Append(et)
	mb.importTrackItem = gtk.NewMenuItemWithLabel("Import CSV Track")
	mb.importTrackItem.Connect("activate", importTrackCB)
	trackMenu.Append(mb.importTrackItem)
	st := gtk.NewMenuItemWithLabel("Export Track as PNG")
	st.Connect("activate", exportTrackImageCB)
	trackMenu.Append(st)

	trackMenu.Append(gtk.NewSeparatorMenuItem())

	mb.trackShowDrone = gtk.NewCheckMenuItemWithLabel("Show Drone Positions")
	mb.trackShowDrone.SetActive(true)
	trackMenu.Append(mb.trackShowDrone)
	mb.trackShowDrone.Connect("activate", func() {
		trackChart.showDrone = mb.trackShowDrone.GetActive()
		trackChart.drawTrack()
	})
	mb.trackShowPath = gtk.NewCheckMenuItemWithLabel("Show Drone Path")
	mb.trackShowPath.SetActive(true)
	trackMenu.Append(mb.trackShowPath)
	mb.trackShowPath.Connect("activate", func() {
		trackChart.showPath = mb.trackShowPath.GetActive()
		trackChart.drawTrack()
	})

	videoItem := gtk.NewMenuItemWithLabel("Video")
	mb.Append(videoItem)
	videoMenu := gtk.NewMenu()
	videoItem.SetSubmenu(videoMenu)

	mb.recVidItem = gtk.NewMenuItemWithLabel("Record Video")
	mb.recVidItem.Connect("activate", recordVideoCB)
	videoMenu.Append(mb.recVidItem)
	mb.stopRecVidItem = gtk.NewMenuItemWithLabel("Stop Recording Video")
	mb.stopRecVidItem.Connect("activate", stopRecordingVideoCB)
	mb.stopRecVidItem.SetSensitive(false)
	videoMenu.Append(mb.stopRecVidItem)

	imagesItem := gtk.NewMenuItemWithLabel("Images")
	mb.Append(imagesItem)
	imagesMenu := gtk.NewMenu()
	imagesItem.SetSubmenu(imagesMenu)

	tp := gtk.NewMenuItemWithLabel("Take Photo")
	tp.Connect("activate", takePhotoCB)
	imagesMenu.Append(tp)
	sp := gtk.NewMenuItemWithLabel("Save Photo(s)")
	sp.Connect("activate", saveAllPhotosCB)
	imagesMenu.Append(sp)

	helpItem := gtk.NewMenuItemWithLabel("Help")
	mb.Append(helpItem)
	helpMenu := gtk.NewMenu()
	helpItem.SetSubmenu(helpMenu)

	oh := gtk.NewMenuItemWithLabel("Online Help")
	oh.Connect("activate", func() { openBrowser(appHelpURL) })
	helpMenu.Append(oh)

	helpMenu.Append(gtk.NewSeparatorMenuItem())

	ab := gtk.NewMenuItemWithLabel("About")
	ab.Connect("activate", aboutCB)
	helpMenu.Append(ab)

	mb.disableFlightMenus()

	return mb
}

func (mb *menuBarT) enableFlightMenus() {
	mb.disconnectItem.SetSensitive(true)
	mb.connectItem.SetSensitive(false)
	mb.flightItem.SetSensitive(true)
	mb.importTrackItem.SetSensitive(false)
}

func (mb *menuBarT) disableFlightMenus() {
	mb.disconnectItem.SetSensitive(false)
	mb.connectItem.SetSensitive(true)
	mb.flightItem.SetSensitive(false)
	mb.importTrackItem.SetSensitive(true)
}
