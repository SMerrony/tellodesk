package main

import (
	"time"
)

func (app *tdApp) connectCB(s string, i interface{}) {
	err := drone.ControlConnectDefault()
	if err != nil {
		alertDialog(
			app.mainPanel,
			errorSev,
			`Could not connect to Drone.

Check that you have a Wifi connection 
to the Tello network.`)
		return // Comment this for GUI testing
	}

	app.startVideo()

	stickChan, _ = drone.StartStickListener()
	go readJoystick(false, jsStopChan) // FIXME - if defined & opened ok!

	app.trackChart.track = newTrack()
	app.liveTrackerTimer = app.Gui().TimerManager.SetInterval(500*time.Millisecond, true, app.liveTrackerTCB)

	fdChan, _ = drone.StreamFlightData(false, fdPeriodMs)
	go app.fdListener()

	app.enableFlightMenus()
	app.statusBar.connectionLab.SetText(" Connected ")
}

func (app *tdApp) disconnectCB(s string, i interface{}) {
	drone.VideoDisconnect()
	drone.ControlDisconnect()
	app.Gui().TimerManager.ClearTimeout(app.liveTrackerTimer)
	jsStopChan <- true         // stop the joystick listener goroutine
	fdStopChan <- true         // stop the flight data listener goroutine
	vrStopChan <- true         // stop the video restarter goroutine
	app.stopNewPicChan <- true // stop the video image updater goroutine
	app.disableFlightMenus()
	app.statusBar.connectionLab.SetText(" Disconnected ")
	app.buildFeed()
	app.feedTab.SetContent(app.feed)
}
