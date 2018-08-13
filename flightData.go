package main

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive
// It is started by connectCB() in droneCBs.go when the Tello is connected
func (app *tdApp) fdListener() {
	for {
		tmpFd := <-fdChan
		app.flightDataMu.Lock()
		app.flightData = tmpFd
		app.flightDataMu.Unlock()
		currentTrack.addPositionIfChanged(tmpFd)
		//app.Dispatch("fdUpdate", nil)
	}
}
