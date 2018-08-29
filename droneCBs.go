package main

import "time"

var feedUpdateTimer int

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
	app.disconnectItem.SetEnabled(true)
	app.connectItem.SetEnabled(false)

	app.startVideo()

	feedUpdateTimer = app.Gui().TimerManager.SetInterval(33*time.Millisecond, true, app.updateFeedTCB)
	//app.Subscribe("feedUpdate", app.feedUpdateCB)

	stickChan, _ = drone.StartStickListener()
	go readJoystick(false, jsStopChan) // FIXME - if defined & opened ok!

	app.importTrackItem.SetEnabled(false)
	app.trackChart.track = newTrack()

	fdChan, _ = drone.StreamFlightData(false, fdPeriodMs)
	go app.fdListener()

	app.statusBar.connectionLab.SetText(" Connected ")
}

func (app *tdApp) disconnectCB(s string, i interface{}) {
	drone.ControlDisconnect()
	app.disconnectItem.SetEnabled(false)
	app.connectItem.SetEnabled(true)
	app.importTrackItem.SetEnabled(true)
	jsStopChan <- true // stop the joystick listener goroutine
	fdStopChan <- true // stop the flight data listener goroutine
	vrStopChan <- true // stop the video restarter goroutine
	app.Gui().TimerManager.ClearTimeout(feedUpdateTimer)
	app.statusBar.connectionLab.SetText(" Disconnected ")
}
