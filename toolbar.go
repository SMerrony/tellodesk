package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

// toolBarT also holds a single message label for urgent notifications to appear at
// the top of the screen.
type toolBarT struct {
	*gtk.Toolbar
	goHomeBtn    *gtk.ToolButton
	messageLabel *gtk.Label
}

func buildToolBar() (tb *toolBarT) {
	tb = new(toolBarT)
	tb.Toolbar = gtk.NewToolbar()
	tb.SetStyle(gtk.TOOLBAR_ICONS)

	stopBtn := gtk.NewToolButtonFromStock(gtk.STOCK_MEDIA_PAUSE)
	stopBtn.SetLabel("Hover")
	stopBtn.Connect("clicked", func() { drone.Hover() })
	tb.Add(stopBtn)

	setHomeBtn := gtk.NewToolButtonFromStock(gtk.STOCK_HOME)
	setHomeBtn.SetLabel("Set Home")
	setHomeBtn.Connect("clicked", func() {
		drone.SetHome()
		tb.goHomeBtn.SetSensitive(true)
	})
	tb.Add(setHomeBtn)

	tb.goHomeBtn = gtk.NewToolButtonFromStock(gtk.STOCK_GO_BACK)
	tb.goHomeBtn.SetLabel("Return Home")
	tb.goHomeBtn.Connect("clicked", func() { drone.AutoFlyToXY(0, 0) })
	tb.goHomeBtn.SetSensitive(false)
	tb.Add(tb.goHomeBtn)

	tb.Add(gtk.NewSeparatorToolItem())

	tb.messageLabel = gtk.NewLabel("(No message)")
	tb.messageLabel.SetWidthChars(40)

	mli := gtk.NewToolItem()
	mli.Add(tb.messageLabel)
	mli.SetBorderWidth(1)
	tb.Add(mli)

	return tb
}

func (tb *toolBarT) clearMessage() {
	tb.messageLabel.SetLabel("")
}

func (tb *toolBarT) setMessage(msg string, severity severityType) {
	tb.messageLabel.SetLabel(msg)
	switch severity {
	case infoSev:
		tb.messageLabel.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))
		tb.messageLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("black"))
	case warningSev:
		tb.messageLabel.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("yellow"))
		tb.messageLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("black"))
	case errorSev, criticalSev:
		tb.messageLabel.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("red"))
		tb.messageLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	}
	// tb.messageLabel.SetWidth(msgBoxWidth)
}
