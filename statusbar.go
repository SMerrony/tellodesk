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
	connectionLab, heightLab, batteryPctLab, wifiStrLab, photosLab *gtk.Label
}

func buildStatusbar() (sb *statusBarT) {
	sb = new(statusBarT)
	sb.HBox = gtk.NewHBox(true, 2)

	clf := gtk.NewFrame("")
	sb.connectionLab = gtk.NewLabel("Disconnected") //NewFixedLabel(" Disconnected ", color.RGBA{255, 255, 255, 255})
	clf.Add(sb.connectionLab)
	sb.Add(clf)

	hlf := gtk.NewFrame("")
	sb.heightLab = gtk.NewLabel("Height: 00.0m")
	hlf.Add(sb.heightLab)
	sb.Add(hlf)

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
	flightDataMu.RLock()
	sb.heightLab.SetLabel(fmt.Sprintf("Height: %.1fm", float32(flightData.Height)/10))
	sb.batteryPctLab.SetLabel(fmt.Sprintf("Battery: %d%%", flightData.BatteryPercentage))
	if flightData.BatteryPercentage < 30 {
		sb.batteryPctLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("red"))
	} else {
		sb.batteryPctLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("black"))
	}
	sb.wifiStrLab.SetLabel(fmt.Sprintf("Wifi Strength: %d%%", flightData.WifiStrength))
	if flightData.WifiStrength < 50 {
		sb.wifiStrLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("red"))
	} else {
		sb.wifiStrLab.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("black"))
	}
	flightDataMu.RUnlock()
	sb.photosLab.SetLabel(fmt.Sprintf("Buffered Photos: %d", drone.NumPics()))
}
