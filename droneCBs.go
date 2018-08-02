package main

func (app *tdApp) connectCB(s string, i interface{}) {
	err := drone.ControlConnectDefault()
	if err != nil {
		// 		app.Gui().Add(alertDialog(
		// 			errorSev,
		// 			`Could not connect to Drone.

		// Check that you have a Wifi connection
		// to the Tello network.`))
		var ad dialogWin
		ad.alertDialog(
			app,
			errorSev,
			`Could not connect to Drone.

Check that you have a Wifi connection 
to the Tello network.`)

		return // Comment this for GUI testing
	}
	app.disconnectItem.SetEnabled(true)
	app.connectItem.SetEnabled(false)
}

func (app *tdApp) diconnectCB(s string, i interface{}) {
	drone.ControlDisconnect()
	app.disconnectItem.SetEnabled(false)
	app.connectItem.SetEnabled(true)
}
