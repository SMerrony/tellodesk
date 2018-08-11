package main

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
	stickChan, _ = drone.StartStickListener()
	go readJoystick(false, jsStopChan)
}

func (app *tdApp) diconnectCB(s string, i interface{}) {
	drone.ControlDisconnect()
	app.disconnectItem.SetEnabled(false)
	app.connectItem.SetEnabled(true)
	jsStopChan <- true // stop the joystick listener goroutine
}
