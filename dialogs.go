/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"github.com/mattn/go-gtk/gtk"
)

func messageDialog(win *gtk.Window, sev gtk.MessageType, msg string) {
	alert := gtk.NewMessageDialog(
		win,
		gtk.DIALOG_MODAL+gtk.DIALOG_DESTROY_WITH_PARENT,
		sev,
		gtk.BUTTONS_CLOSE,
		msg)
	alert.SetTitle(appName)
	alert.SetIcon(iconPixbuf)
	alert.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	alert.Run()
	alert.Destroy()
}
