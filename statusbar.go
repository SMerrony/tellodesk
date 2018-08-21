package main

import (
	"fmt"
	"image/color"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

type statusbar struct {
	*gui.Panel
	connectionLab, heightLab, batteryPctLab, wifiStrLab *FixedLabel
}

func buildStatusbar(parent *gui.Panel) (sb *statusbar) {
	sb = new(statusbar)
	sb.Panel = gui.NewPanel(parent.Width(), 30)

	hbl := gui.NewHBoxLayout()
	hbl.SetSpacing(4)
	sb.SetLayout(hbl)

	labStyle := gui.PanelStyle{
		Margin:      gui.RectBounds{Top: 3, Right: 3, Bottom: 3, Left: 3},
		Border:      gui.RectBounds{Top: 1, Right: 1, Bottom: 1, Left: 1},
		Padding:     gui.RectBounds{Top: 2, Right: 2, Bottom: 2, Left: 2},
		BorderColor: math32.Color4Name("black"),
		BgColor:     math32.Color4Name("dark gray"),
	}
	padCol := math32.ColorName("gray")
	params := gui.HBoxLayoutParams{Expand: 0}
	padParams := gui.HBoxLayoutParams{Expand: 1}

	// sb.connectionLab = gui.NewLabel(" Connection Status ")
	// sb.connectionLab.ApplyStyle(&labStyle)
	// sb.connectionLab.SetColor(math32.NewColor("white"))
	// sb.connectionLab.SetPaddingsColor(&padCol)
	// sb.connectionLab.SetLayoutParams(&params)
	sb.connectionLab = NewFixedLabel(" Disconnected ", color.RGBA{255, 255, 255, 255})
	sb.connectionLab.ApplyStyle(&labStyle)
	sb.connectionLab.SetPaddingsColor(&padCol)
	sb.connectionLab.SetLayoutParams(&params)
	sb.Add(sb.connectionLab)

	padder := gui.NewLabel("")
	padder.SetLayoutParams(&padParams)
	sb.Add(padder)

	sb.heightLab = NewFixedLabel(" Height: 00.0m ", color.RGBA{255, 255, 255, 255})
	sb.heightLab.ApplyStyle(&labStyle)
	sb.heightLab.SetPaddingsColor(&padCol)
	sb.heightLab.SetLayoutParams(&params)
	sb.Add(sb.heightLab)

	padder2 := gui.NewLabel("")
	padder2.SetLayoutParams(&padParams)
	sb.Add(padder2)

	sb.batteryPctLab = NewFixedLabel(" Battery: 000% ", color.RGBA{255, 255, 255, 255})
	sb.batteryPctLab.ApplyStyle(&labStyle)
	sb.batteryPctLab.SetPaddingsColor(&padCol)
	sb.batteryPctLab.SetLayoutParams(&params)
	sb.Add(sb.batteryPctLab)

	padder3 := gui.NewLabel("")
	padder3.SetLayoutParams(&padParams)
	sb.Add(padder3)

	sb.wifiStrLab = NewFixedLabel(" Wifi Strength: 000% ", color.RGBA{255, 255, 255, 255})
	sb.wifiStrLab.ApplyStyle(&labStyle)
	sb.wifiStrLab.SetPaddingsColor(&padCol)
	sb.wifiStrLab.SetLayoutParams(&params)
	sb.Add(sb.wifiStrLab)

	return sb
}

// func (app *tdApp) updateStatusBar(s string, ev interface{}) {
// 	app.flightDataMu.RLock()
// 	app.statusBar.heightLab.SetText(fmt.Sprintf(" Height: %.1fm ", float32(app.flightData.Height)/10))
// 	app.statusBar.batteryPctLab.SetText(fmt.Sprintf(" Battery: %d%% ", app.flightData.BatteryPercentage))
// 	app.statusBar.wifiStrLab.SetText(fmt.Sprintf(" Wifi Strength: %d%% ", app.flightData.WifiStrength))
// 	app.flightDataMu.RUnlock()
// 	app.statusBar.SetChanged(true)
// }

func (app *tdApp) updateStatusBarTCB(cb interface{}) {
	app.flightDataMu.RLock()
	app.statusBar.heightLab.SetText(fmt.Sprintf(" Height: %.1fm ", float32(app.flightData.Height)/10))
	app.statusBar.batteryPctLab.SetText(fmt.Sprintf(" Battery: %d%% ", app.flightData.BatteryPercentage))
	app.statusBar.wifiStrLab.SetText(fmt.Sprintf(" Wifi Strength: %d%% ", app.flightData.WifiStrength))
	app.flightDataMu.RUnlock()
	//app.statusBar.SetChanged(true)
}
