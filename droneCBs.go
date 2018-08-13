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
	app.disconnectItem.SetEnabled(true)
	app.connectItem.SetEnabled(false)

	app.startVideo()

	app.Gui().TimerManager.SetInterval(33*time.Millisecond, true, app.updateFeedTCB)

	stickChan, _ = drone.StartStickListener()
	go readJoystick(false, jsStopChan) // FIXME - if defined & opened ok!

	currentTrack = newTrack()

	fdChan, _ = drone.StreamFlightData(false, fdPeriodMs)
	go app.fdListener()
}

func (app *tdApp) diconnectCB(s string, i interface{}) {
	drone.ControlDisconnect()
	app.disconnectItem.SetEnabled(false)
	app.connectItem.SetEnabled(true)
	jsStopChan <- true // stop the joystick listener goroutine
}
