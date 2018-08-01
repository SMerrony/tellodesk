package main

func (app *tdApp) connectCB(s string, i interface{}) {
	err := drone.ControlConnectDefault()
	if err != nil {
		app.Gui().Add(runAlert(errorSev, "Could not connect to Drone"))
	}
}
