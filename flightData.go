package main

import "fmt"

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive
// It is started by connectCB() in droneCBs.go when the Tello is connected
func (app *tdApp) fdListener() {
	for {
		fd := <-fdChan
		currentTrack.addPositionIfChanged(fd)
		app.statusBar.heightLab.SetText(fmt.Sprintf(" Height: %.1fm ", float32(fd.Height)/10))
		app.statusBar.batteryPctLab.SetText(fmt.Sprintf(" Battery: %d%% ", fd.BatteryPercentage))
		app.statusBar.wifiStrLab.SetText(fmt.Sprintf(" Wifi Strength: %d%% ", fd.WifiStrength))
		app.statusBar.SetChanged(true)
	}
}
