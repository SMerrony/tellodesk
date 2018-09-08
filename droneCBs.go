package main

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func connectCB() {

	err := drone.ControlConnectDefault()
	if err != nil {
		alert := gtk.NewMessageDialog(win, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE,
			`Could not connect to Drone.

Check that you have a Wifi connection 
to the Tello network.`)
		alert.SetTitle(appName)
		alert.Run()
		alert.Destroy()
		return // Comment this for GUI testing
	}

	startVideo()

	stickChan, _ = drone.StartStickListener()
	go readJoystick(false, jsStopChan) // FIXME - if defined & opened ok!

	trackChart.track = newTrack()
	glib.TimeoutAdd(500, liveTrackerTCB) // start the live tracker, cancelled via liveTrackStopChan

	fdChan, _ = drone.StreamFlightData(false, fdPeriodMs)
	go fdListener()

	menuBar.enableFlightMenus()
	statusBar.connectionLab.SetText("Connected")
}

func disconnectCB() {
	drone.VideoDisconnect()
	drone.ControlDisconnect()

	select {
	case jsStopChan <- true: // stop the joystick listener goroutine
	default:
	}
	select {
	case fdStopChan <- true: // stop the flight data listener goroutine
	default:
	}
	select {
	case vrStopChan <- true: // stop the video restarter goroutine
	default:
	}
	select {
	case liveTrackStopChan <- true: // stop the live tracker
	default:
	}
	select {
	case stopFeedImageChan <- true: // stop the video image updater goroutine
	default:
	}

	menuBar.disableFlightMenus()
	statusBar.connectionLab.SetText(" Disconnected ")
}
