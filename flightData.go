package main

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive.
// It is started by connectCB() in droneCBs.go when the Tello is connected.
func (app *tdApp) fdListener() {
	for {
		select {
		case tmpFd := <-fdChan:
			app.flightDataMu.Lock()
			app.flightData = tmpFd
			app.flightDataMu.Unlock()
			if tmpFd.DownVisualState {
				app.Log().Info("Down visual state")
			}
			if tmpFd.OnGround {
				app.Log().Info("On Ground")
			}
			if tmpFd.LightStrength == 0 {
				app.trackChart.track.addPositionIfChanged(tmpFd)
			}
		case <-fdStopChan:
			return
		}
	}
}

// updateMessage should be run periodically to check for condition we should alert the user about.
func (app *tdApp) updateMessage(cb interface{}) {
	var (
		msg string
		sev severityType
	)
	app.flightDataMu.RLock()
	// is order of priority, descending...
	switch {
	case app.flightData.OverTemp:
		msg = "Maximum Temperature Exceeded"
		sev = criticalSev
	case app.flightData.BatteryCritical:
		msg = "Battery Critical"
		sev = criticalSev
	case app.flightData.WifiStrength < 30:
		msg = "Wifi Strength Below 30%"
		sev = criticalSev
	case app.flightData.BatteryLow:
		msg = "Battery Low"
		sev = warningSev
	case app.flightData.WifiStrength < 50:
		msg = "Wifi Strength Below 50%"
		sev = infoSev
	case app.flightData.LightStrength == 1:
		msg = "Low Light"
		sev = infoSev
	}
	app.flightDataMu.RUnlock()
	if msg == "" {
		app.toolBar.clearMessage()
	} else {
		app.toolBar.setMessage(msg, sev)
	}
}
