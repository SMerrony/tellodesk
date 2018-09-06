package main

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/math32"
)

// toolbarT also holds a single message label for urgent notifications to appear at
// the top of the screen.
type toolbarT struct {
	*gui.Panel
	messageLabel *gui.ImageLabel
}

const msgBoxWidth = 350.0

func (app *tdApp) buildToolbar() (tb *toolbarT) {
	tb = new(toolbarT)
	tb.Panel = gui.NewPanel(videoWidth, 28)
	tb.SetContentWidth(videoWidth)
	//tb.SetBorders(1, 1, 1, 1)
	tb.SetMargins(1, 1, 1, 1)

	hbl := gui.NewHBoxLayout()
	hbl.SetSpacing(4)
	tb.SetLayout(hbl)

	stopBtn := gui.NewButton("")
	stopBtn.SetIcon(icon.Pause)
	stopBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		drone.Hover()
	})
	tb.Add(stopBtn)

	cameraBtn := gui.NewButton("")
	cameraBtn.SetIcon(icon.CameraAlt)
	cameraBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		drone.TakePicture()
	})
	tb.Add(cameraBtn)

	setHomeBtn := gui.NewButton("")
	setHomeBtn.SetIcon(icon.AddLocation)
	setHomeBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		drone.SetHome()
	})
	tb.Add(setHomeBtn)

	goHomeBtn := gui.NewButton("")
	goHomeBtn.SetIcon(icon.Place)
	goHomeBtn.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		if drone.IsHomeSet() {
			drone.AutoFlyToXY(0, 0)
		} else {
			alertDialog(app.mainPanel, warningSev, "Home position not set")
		}
	})
	tb.Add(goHomeBtn)

	spacer := gui.NewPanel(1, 1)
	spacer.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 2})
	tb.Add(spacer)

	tb.messageLabel = gui.NewImageLabel("")
	tb.messageLabel.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	tb.messageLabel.SetWidth(msgBoxWidth)
	tb.messageLabel.SetBorders(2, 2, 2, 2)
	tb.clearMessage()

	tb.Add(tb.messageLabel)

	return tb
}

func (tb *toolbarT) clearMessage() {
	tb.messageLabel.SetText("")
	tb.messageLabel.SetBgColor(math32.NewColor("white"))
}

func (tb *toolbarT) setMessage(msg string, severity severityType) {
	tb.messageLabel.SetText(msg)
	switch severity {
	case infoSev:
		tb.messageLabel.SetBgColor(math32.NewColor("white"))
		tb.messageLabel.SetColor(math32.NewColor("black"))
	case warningSev:
		tb.messageLabel.SetBgColor(math32.NewColor("yellow"))
		tb.messageLabel.SetColor(math32.NewColor("black"))
	case errorSev, criticalSev:
		tb.messageLabel.SetBgColor(math32.NewColor("red"))
		tb.messageLabel.SetColor(math32.NewColor("white"))
	}
	tb.messageLabel.SetWidth(msgBoxWidth)
}
