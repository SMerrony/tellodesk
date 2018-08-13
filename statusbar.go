package main

import (
	"fmt"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

type statusbar struct {
	*gui.Panel
	connectionLab, heightLab,
	batteryPctLab, wifiStrLab *gui.Label
}

func buildStatusbar(parent *gui.Panel) (sb *statusbar) {
	sb = new(statusbar)
	sb.Panel = gui.NewPanel(parent.Width(), 30)

	hbl := gui.NewHBoxLayout()
	hbl.SetSpacing(4)
	sb.SetLayout(hbl)

	labStyle := gui.PanelStyle{
		Margin:      gui.RectBounds{3, 3, 3, 3},
		Border:      gui.RectBounds{1, 1, 1, 1},
		Padding:     gui.RectBounds{2, 2, 2, 2},
		BorderColor: math32.Color4Name("black"),
		BgColor:     math32.Color4Name("dark gray"),
	}
	padCol := math32.ColorName("gray")
	params := gui.HBoxLayoutParams{Expand: 0}
	padParams := gui.HBoxLayoutParams{Expand: 1}

	sb.connectionLab = gui.NewLabel(" Connection Status ")
	sb.connectionLab.ApplyStyle(&labStyle)
	sb.connectionLab.SetColor(math32.NewColor("white"))
	sb.connectionLab.SetPaddingsColor(&padCol)
	sb.connectionLab.SetLayoutParams(&params)
	sb.Add(sb.connectionLab)

	padder := gui.NewLabel("")
	padder.SetLayoutParams(&padParams)
	sb.Add(padder)

	sb.heightLab = gui.NewLabel(" Height: 00.0 m ")
	sb.heightLab.ApplyStyle(&labStyle)
	sb.heightLab.SetPaddingsColor(&padCol)
	sb.heightLab.SetLayoutParams(&params)
	sb.Add(sb.heightLab)

	padder2 := gui.NewLabel("")
	padder2.SetLayoutParams(&padParams)
	sb.Add(padder2)

	sb.batteryPctLab = gui.NewLabel(" Battery: 000% ")
	sb.batteryPctLab.ApplyStyle(&labStyle)
	sb.batteryPctLab.SetPaddingsColor(&padCol)
	sb.batteryPctLab.SetLayoutParams(&params)
	sb.Add(sb.batteryPctLab)

	padder3 := gui.NewLabel("")
	padder3.SetLayoutParams(&padParams)
	sb.Add(padder3)

	sb.wifiStrLab = gui.NewLabel(" Wifi Strength: 000% ")
	sb.wifiStrLab.ApplyStyle(&labStyle)
	sb.wifiStrLab.SetPaddingsColor(&padCol)
	sb.wifiStrLab.SetLayoutParams(&params)
	sb.Add(sb.wifiStrLab)

	return sb
}

func (app *tdApp) updateStatusBar(s string, ev interface{}) {
	flightDataMu.RLock()
	app.statusBar.heightLab.SetText(fmt.Sprintf(" Height: %.1fm ", float32(flightData.Height)/10))
	app.statusBar.batteryPctLab.SetText(fmt.Sprintf(" Battery: %d%% ", flightData.BatteryPercentage))
	app.statusBar.wifiStrLab.SetText(fmt.Sprintf(" Wifi Strength: %d%% ", flightData.WifiStrength))
	flightDataMu.RUnlock()
	app.statusBar.SetChanged(true)
}
