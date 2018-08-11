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
	sb.Panel = gui.NewPanel(parent.Width(), parent.Height())

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

	sb.connectionLab = gui.NewLabel("Connection Status")
	sb.connectionLab.ApplyStyle(&labStyle)
	sb.connectionLab.SetColor(math32.NewColor("white"))
	sb.connectionLab.SetPaddingsColor(&padCol)
	sb.connectionLab.SetLayoutParams(&params)
	sb.Add(sb.connectionLab)

	padder := gui.NewLabel("")
	padder.SetLayoutParams(&padParams)
	sb.Add(padder)

	sb.heightLab = gui.NewLabel("Height:   . m")
	sb.heightLab.ApplyStyle(&labStyle)
	sb.heightLab.SetPaddingsColor(&padCol)
	sb.heightLab.SetLayoutParams(&params)
	sb.Add(sb.heightLab)

	padder2 := gui.NewLabel("")
	padder2.SetLayoutParams(&padParams)
	sb.Add(padder2)

	sb.batteryPctLab = gui.NewLabel("Battery:   %")
	sb.batteryPctLab.ApplyStyle(&labStyle)
	sb.batteryPctLab.SetPaddingsColor(&padCol)
	sb.batteryPctLab.SetLayoutParams(&params)
	sb.Add(sb.batteryPctLab)

	padder3 := gui.NewLabel("")
	padder3.SetLayoutParams(&padParams)
	sb.Add(padder3)

	sb.wifiStrLab = gui.NewLabel("Wifi Strength:   %")
	sb.wifiStrLab.ApplyStyle(&labStyle)
	sb.wifiStrLab.SetPaddingsColor(&padCol)
	sb.wifiStrLab.SetLayoutParams(&params)
	sb.Add(sb.wifiStrLab)

	return sb
}

// fdListener should be run as a Goroutine to consume FD updates on the chan as they arrive
func (app *tdApp) fdListener() {
	for {
		fd := <-fdChan
		app.statusBar.heightLab.SetText(fmt.Sprintf("Height: %.1fm", float32(fd.Height)/10))
		app.statusBar.batteryPctLab.SetText(fmt.Sprintf("Battery: %d%%", fd.BatteryPercentage))
		app.statusBar.wifiStrLab.SetText(fmt.Sprintf("Wifi Strength: %d%%", fd.WifiStrength))
		app.statusBar.SetChanged(true)
	}
}
