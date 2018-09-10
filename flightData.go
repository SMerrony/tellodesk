package main

import "log"

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive.
// It is started by connectCB() in droneCBs.go when the Tello is connected.
func fdListener() {
	for {
		select {
		case tmpFd := <-fdChan:
			flightDataMu.Lock()
			flightData = tmpFd
			flightDataMu.Unlock()
			if tmpFd.DownVisualState {
				log.Println("Down visual state")
			}
			if tmpFd.OnGround {
				log.Println("On Ground")
			}
			if tmpFd.LightStrength == 0 {
				trackChart.track.addPositionIfChanged(tmpFd)
			}
		case <-fdStopChan:
			return
		}
	}
}

// // updateMessage should be run periodically to check for condition we should alert the user about.
// func updateMessageCB() bool {
// 	var (
// 		msg string
// 		sev severityType
// 	)
// 	flightDataMu.RLock()
// 	// is order of priority, descending...
// 	switch {
// 	case flightData.OverTemp:
// 		msg = "Maximum Temperature Exceeded"
// 		sev = criticalSev
// 	case flightData.BatteryCritical:
// 		msg = "Battery Critical"
// 		sev = criticalSev
// 	case flightData.WifiStrength < 30:
// 		msg = "Wifi Strength Below 30%"
// 		sev = criticalSev
// 	case flightData.BatteryLow:
// 		msg = "Battery Low"
// 		sev = warningSev
// 	case flightData.WifiStrength < 50:
// 		msg = "Wifi Strength Below 50%"
// 		sev = infoSev
// 	case flightData.LightStrength == 1:
// 		msg = "Low Light"
// 		sev = infoSev
// 	}
// 	flightDataMu.RUnlock()
// 	if msg == "" {
// 		//toolBar.clearMessage()
// 	} else {
// 		//toolBar.setMessage(msg, sev)
// 		_ = sev
// 	}
// 	return true // continue the timer
// }
