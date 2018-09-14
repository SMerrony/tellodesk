/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"github.com/mattn/go-gtk/gtk"
)

type menuBarT struct {
	*gtk.MenuBar
	connectItem, disconnectItem             *gtk.MenuItem
	navItem, goHomeItem, flightItem         *gtk.MenuItem
	sportsModeItem                          *gtk.CheckMenuItem
	importTrackItem                         *gtk.MenuItem
	imagingItem, recVidItem, stopRecVidItem *gtk.MenuItem
	trackShowDrone, trackShowPath           *gtk.CheckMenuItem
}

func buildMenu() (mb *menuBarT) {

	mb = new(menuBarT)
	mb.MenuBar = gtk.NewMenuBar()

	// File

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

	// Drone

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

	mb.sportsModeItem = gtk.NewCheckMenuItemWithLabel("Sports (Fast) Mode")
	mb.sportsModeItem.Connect("activate", toggleSportsModeCB)
	flightMenu.Append(mb.sportsModeItem)

	// Navigation

	mb.navItem = gtk.NewMenuItemWithLabel("Navigation")
	mb.Append(mb.navItem)
	navMenu := gtk.NewMenu()
	mb.navItem.SetSubmenu(navMenu)

	sh := gtk.NewMenuItemWithLabel("Set Home Position")
	sh.Connect("activate", func() {
		drone.SetHome()
		mb.goHomeItem.SetSensitive(true)
	})
	navMenu.Append(sh)
	mb.goHomeItem = gtk.NewMenuItemWithLabel("Return to Home")
	mb.goHomeItem.SetSensitive(false)
	mb.goHomeItem.Connect("activate", func() { drone.AutoFlyToXY(0, 0) })
	navMenu.Append(mb.goHomeItem)

	// Track

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

	// Imaging

	mb.imagingItem = gtk.NewMenuItemWithLabel("Imaging")
	mb.Append(mb.imagingItem)
	imagingMenu := gtk.NewMenu()
	mb.imagingItem.SetSubmenu(imagingMenu)

	mb.recVidItem = gtk.NewMenuItemWithLabel("Record Video")
	mb.recVidItem.Connect("activate", recordVideoCB)
	imagingMenu.Append(mb.recVidItem)
	mb.stopRecVidItem = gtk.NewMenuItemWithLabel("Stop Recording Video")
	mb.stopRecVidItem.Connect("activate", stopRecordingVideoCB)
	mb.stopRecVidItem.SetSensitive(false)
	imagingMenu.Append(mb.stopRecVidItem)

	imagingMenu.Append(gtk.NewSeparatorMenuItem())

	tp := gtk.NewMenuItemWithLabel("Take Photo")
	tp.Connect("activate", takePhotoCB)
	imagingMenu.Append(tp)
	sp := gtk.NewMenuItemWithLabel("Save Photo(s)")
	sp.Connect("activate", saveAllPhotosCB)
	imagingMenu.Append(sp)

	// Help

	helpItem := gtk.NewMenuItemWithLabel("Help")
	mb.Append(helpItem)
	helpMenu := gtk.NewMenu()
	helpItem.SetSubmenu(helpMenu)

	jh := gtk.NewMenuItemWithLabel("Joystick Functions")
	jh.Connect("activate", func() { joystickHelpCB() })
	helpMenu.Append(jh)

	oh := gtk.NewMenuItemWithLabel("Online Help")
	oh.Connect("activate", func() { openBrowser(appHelpURL) })
	helpMenu.Append(oh)

	helpMenu.Append(gtk.NewSeparatorMenuItem())

	ab := gtk.NewMenuItemWithLabel("About")
	ab.Connect("activate", aboutCB)
	helpMenu.Append(ab)

	mb.disableFlightMenus() // At startup all flight functions are disabled

	return mb
}

func (mb *menuBarT) enableFlightMenus() {
	mb.disconnectItem.SetSensitive(true)
	mb.connectItem.SetSensitive(false)
	mb.flightItem.SetSensitive(true)
	mb.navItem.SetSensitive(true)
	mb.imagingItem.SetSensitive(true)
	mb.importTrackItem.SetSensitive(false)
}

func (mb *menuBarT) disableFlightMenus() {
	mb.disconnectItem.SetSensitive(false)
	mb.connectItem.SetSensitive(true)
	mb.flightItem.SetSensitive(false)
	mb.navItem.SetSensitive(false)
	mb.imagingItem.SetSensitive(false)
	mb.importTrackItem.SetSensitive(true)
}
