package main

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive
// It is started by connectCB() in droneCBs.go when the Tello is connected
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
