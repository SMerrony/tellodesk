/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/SMerrony/tello"
)

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive.
// It is started by connectCB() in droneCBs.go when the Tello is connected.
func fdListener(fdChan <-chan tello.FlightData) {
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
				liveTrack.addPositionIfChanged(tmpFd)
			}
		case <-fdStopChan:
			return
		}
	}
}

// updateFlightDataTCB should be run periodically to check for condition we should alert the user about
// and update the flight status display.
func updateFlightDataTCB() bool {
	msg := ""
	flightDataMu.RLock()

	// first, the message overlaid on the video display

	if len(flightData.SSID) > 0 {
		switch {
		case flightData.BatteryCritical:
			msg = "Battery Critical  "
		case flightData.BatteryLow:
			msg = "Battery Low  "
		}
		if flightData.WifiStrength < 30 {
			msg += "Wifi < 30%  "
		} else if flightData.WifiStrength < 50 {
			msg += "Wifi < 50%  "
		}
		if flightData.LightStrength == 1 {
			msg += "Low Light  "
		}
		if flightData.WindState {
			msg += "Too Windy  "
		}
	}

	// now the flight status display
	statFields[fYaw].value.SetText(fmt.Sprintf("%d°", flightData.IMU.Yaw))
	statFields[fLoBattThres].value.SetText(fmt.Sprintf("%d%%", flightData.LowBatteryThreshold))
	statFields[fTemp].value.SetText(fmt.Sprintf("%dC", flightData.IMU.Temperature))
	statFields[fDrvdSpeed].value.SetText(fmt.Sprintf("%.1fm/s", math.Sqrt(float64(flightData.NorthSpeed*flightData.NorthSpeed)+float64(flightData.EastSpeed*flightData.EastSpeed))))
	statFields[fVertSpeed].value.SetText(fmt.Sprintf("%dm/s", flightData.VerticalSpeed))
	statFields[fGndSpeed].value.SetText(fmt.Sprintf("%dm/s", flightData.GroundSpeed))
	statFields[fFwdSpeed].value.SetText(fmt.Sprintf("%dm/s", flightData.NorthSpeed))
	statFields[fLatSpeed].value.SetText(fmt.Sprintf("%dm/s", flightData.EastSpeed))
	statFields[fBattLow].value.SetText(boolToYN(flightData.BatteryLow))
	statFields[fBattCrit].value.SetText(boolToYN(flightData.BatteryCritical))
	statFields[fBattErr].value.SetText(boolToYN(flightData.BatteryState))
	statFields[fXPos].value.SetText(fmt.Sprintf("%.2fm", flightData.MVO.PositionX))
	statFields[fYPos].value.SetText(fmt.Sprintf("%.2fm", flightData.MVO.PositionY))
	statFields[fZPos].value.SetText(fmt.Sprintf("%.2fm", flightData.MVO.PositionZ))
	statFields[fOnGround].value.SetText(boolToYN(flightData.OnGround))
	statFields[fFlying].value.SetText(boolToYN(flightData.Flying))
	statFields[fWindy].value.SetText(boolToYN(flightData.WindState))

	flightDataMu.RUnlock()
	if msg == "" {
		videoWgt.clearMessage()
	} else {
		videoWgt.setMessage(msg)
	}
	return true // continue the timer
}

func boolToYN(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}
