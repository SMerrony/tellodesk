/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"fmt"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

type statusBarT struct {
	*gtk.HBox
	connectionLab, heightLab, homeLab, batteryPctLab, wifiStrLab, photosLab *gtk.Label
}

func buildStatusbar() (sb *statusBarT) {
	sb = new(statusBarT)
	sb.HBox = gtk.NewHBox(false, 2)

	clf := gtk.NewFrame("")
	sb.connectionLab = gtk.NewLabel("Disconnected") //NewFixedLabel(" Disconnected ", color.RGBA{255, 255, 255, 255})
	clf.Add(sb.connectionLab)
	sb.Add(clf)

	htlf := gtk.NewFrame("")
	sb.heightLab = gtk.NewLabel("Height: 00.0m")
	htlf.Add(sb.heightLab)
	sb.Add(htlf)

	holf := gtk.NewFrame("")
	sb.homeLab = gtk.NewLabel("Home: Unset")
	holf.Add(sb.homeLab)
	sb.Add(holf)

	blf := gtk.NewFrame("")
	sb.batteryPctLab = gtk.NewLabel("Battery: 000%")
	sb.batteryPctLab.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("yellow"))
	blf.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("red"))
	blf.Add(sb.batteryPctLab)
	sb.Add(blf)

	wlf := gtk.NewFrame("")
	sb.wifiStrLab = gtk.NewLabel("Wifi Strength: 000%")
	wlf.Add(sb.wifiStrLab)
	sb.Add(wlf)

	plf := gtk.NewFrame("")
	sb.photosLab = gtk.NewLabel("Buffered Photos: 00")
	plf.Add(sb.photosLab)
	sb.Add(plf)

	return sb
}

func (sb *statusBarT) updateStatusBarTCB() {
	if drone.ControlConnected() {
		flightDataMu.RLock()
		if len(flightData.SSID) > 0 {
			sb.connectionLab.SetLabel(fmt.Sprintf("%s - Firmware: %s", flightData.SSID, flightData.Version))
		} else {
			sb.connectionLab.SetLabel("Disconnected")
		}
		sb.heightLab.SetLabel(fmt.Sprintf("Height: %.1fm (Max: %dm)", float32(flightData.Height)/10, flightData.MaxHeight))
		sb.batteryPctLab.SetLabel(fmt.Sprintf("Battery: %d%% (%dmV)", flightData.BatteryPercentage, flightData.BatteryMilliVolts))
		if flightData.BatteryPercentage < 30 {
			sb.batteryPctLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("red"))
		} else {
			sb.batteryPctLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("black"))
		}
		sb.wifiStrLab.SetLabel(fmt.Sprintf("Wifi: %d%% - Interference: %d%%", flightData.WifiStrength, flightData.WifiInterference))
		if flightData.WifiStrength < 50 {
			sb.wifiStrLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("red"))
		} else {
			sb.wifiStrLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("black"))
		}
		flightDataMu.RUnlock()
		// outside the lock...
		if drone.IsHomeSet() {
			if drone.IsAutoXY() {
				sb.homeLab.SetLabel("AUTOFLIGHT")
			} else {
				sb.homeLab.SetLabel("Home: Set")
			}
		} else {
			sb.homeLab.SetLabel("Home: Unset")
		}
	} else {
		// disconnected
		sb.connectionLab.SetLabel("Disconnected")
		sb.heightLab.SetLabel("Height: Unknown")
		sb.batteryPctLab.SetLabel("Battery: Unknown")
		sb.wifiStrLab.SetLabel("Wifi Strength: Unknown")
	}
	sb.photosLab.SetLabel(fmt.Sprintf("Buffered Photos: %d", drone.NumPics()))
}
