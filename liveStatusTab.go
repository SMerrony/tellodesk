/**
 *Copyright (c) 2019 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import "github.com/mattn/go-gtk/gtk"

const (
	rows, cols = 12, 6
)

type liveStatusTabT struct {
	*gtk.Table
}

type sFieldT struct {
	row, col uint
	label    string
	value    *gtk.Label
}

const (
	fYaw = iota
	fLoBattThres
	fTemp
	fDrvdSpeed
	fVertSpeed
	fGndSpeed
	fFwdSpeed
	fLatSpeed
	fBattLow
	fBattCrit
	fBattErr
	fXPos
	fYPos
	fZPos
	fOnGround
	fFlying
	fWindy
	fCount
)

var statFields [fCount]sFieldT

func initFields() {

	statFields[fYaw] = sFieldT{row: 2, col: 0, label: "Yaw:"}
	statFields[fLoBattThres] = sFieldT{row: 2, col: 2, label: "Lo Batt Threshold:"}
	statFields[fTemp] = sFieldT{row: 2, col: 4, label: "Temperature:"}

	statFields[fDrvdSpeed] = sFieldT{row: 4, col: 2, label: "Derived Speed:"}
	statFields[fVertSpeed] = sFieldT{row: 4, col: 4, label: "Vertical Speed:"}

	statFields[fGndSpeed] = sFieldT{row: 5, col: 0, label: "Ground Speed:"}
	statFields[fFwdSpeed] = sFieldT{row: 5, col: 2, label: "Forward Speed:"}
	statFields[fLatSpeed] = sFieldT{row: 5, col: 4, label: "Lateral Speed:"}

	statFields[fBattLow] = sFieldT{row: 7, col: 0, label: "Battery Low:"}
	statFields[fBattCrit] = sFieldT{row: 7, col: 2, label: "Battery Critical:"}
	statFields[fBattErr] = sFieldT{row: 7, col: 4, label: "Battery Error:"}

	statFields[fXPos] = sFieldT{row: 9, col: 0, label: "X Position:"}
	statFields[fYPos] = sFieldT{row: 9, col: 2, label: "Y Position:"}
	statFields[fZPos] = sFieldT{row: 9, col: 4, label: "Z Position:"}

	statFields[fOnGround] = sFieldT{row: 10, col: 0, label: "On Ground:"}
	statFields[fFlying] = sFieldT{row: 10, col: 2, label: "Flying:"}
	statFields[fWindy] = sFieldT{row: 10, col: 4, label: "Windy:"}

	for f := range statFields {
		statFields[f].value = gtk.NewLabel("")
		statFields[f].value.ModifyFontEasy("Sans 16")
		statFields[f].value.SetAlignment(-1, 0.5)
	}
}

func buildLiveStatusTab(w, h int) (st *liveStatusTabT) {
	st = new(liveStatusTabT)
	st.Table = gtk.NewTable(rows, cols, false)
	st.Table.SetRowSpacings(10)
	st.Table.SetColSpacings(10)

	initFields()
	for _, field := range statFields {
		lab := gtk.NewLabel(field.label)
		lab.ModifyFontEasy("Sans 18")
		lab.SetAlignment(1, 0.5)
		st.Table.AttachDefaults(lab, field.col, field.col+1, field.row, field.row+1)
		st.Table.AttachDefaults(field.value, field.col+1, field.col+2, field.row, field.row+1)
	}

	return st
}
