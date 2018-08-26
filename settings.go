package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/g3n/engine/math32"

	"github.com/g3n/engine/gui"
	"gopkg.in/yaml.v2"
)

// settings holds the settings we want to persist across program invocations
type settingsT struct {
	JoystickID   int
	JoystickType string
}

const (
	dialogTitle                           = appName + " Settings"
	settingsWidth, settingsHeight float32 = 500, 200
)

func saveSettings(s settingsT, filename string) error {
	bytes, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}

func loadSettings(filename string) (settingsT, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return settingsT{}, err
	}
	var s settingsT
	err = yaml.Unmarshal(bytes, &s)
	if err != nil {
		return settingsT{}, err
	}
	return s, nil
}

type settingsDlg struct {
	*gui.Window
}

func (app *tdApp) settingsCB(s string, i interface{}) {
	app.settingsDialog()
}

func (app *tdApp) settingsDialog() (win *settingsDlg) {
	win = new(settingsDlg)
	win.Window = gui.NewWindow(settingsWidth, settingsHeight)
	win.SetResizable(false)
	win.SetPaddings(4, 4, 4, 4)
	win.SetTitle(dialogTitle)
	win.SetCloseButton(false)
	win.SetColor(math32.NewColor("Gray"))

	lay := gui.NewGridLayout(3)
	lay.SetAlignH(gui.AlignCenter)
	lay.SetExpandH(true)
	win.SetLayout(lay)

	win.Add(gui.NewLabel(""))
	win.Add(gui.NewLabel("Detected"))
	win.Add(gui.NewLabel("Type"))

	win.Add(gui.NewLabel("Joystick"))
	dDrop := gui.NewDropDown(200, gui.NewImageLabel(""))
	// dDrop.SetMargins(3, 3, 3, 3)
	found := listJoysticks()
	for _, j := range found {
		dDrop.Add(gui.NewImageLabel(j.Name))
	}
	if app.settingsLoaded {
		dDrop.SelectPos(app.settings.JoystickID)
	}
	win.Add(dDrop)

	tDrop := gui.NewDropDown(150, gui.NewImageLabel(""))
	// tDrop.SetMargins(3, 3, 3, 3)
	known := listKnownJoystickTypes()
	for _, k := range known {
		il := gui.NewImageLabel(k.Name)
		tDrop.Add(il)
		if app.settings.JoystickType == k.Name {
			tDrop.SetSelected(il)
		}
	}
	win.Add(tDrop)

	warningLab := gui.NewLabel("You must reconnect to the drone after changing joystick settings")
	warningLab.SetMargins(3, 3, 3, 3)
	warningLab.SetLayoutParams(&gui.GridLayoutParams{ColSpan: 3, AlignH: gui.AlignCenter})
	warningLab.SetBgColor(math32.NewColor("Red"))
	win.Add(warningLab)

	// empty row...
	win.Add(gui.NewLabel(""))
	win.Add(gui.NewLabel(""))
	win.Add(gui.NewLabel(""))

	// buttons...
	win.Add(gui.NewLabel(""))
	cancel := gui.NewButton("Cancel")
	// cancel.SetBorders(1, 1, 1, 1)
	// cancel.SetPaddings(3, 3, 3, 3)
	// cancel.SetMargins(3, 3, 3, 3)
	ok := gui.NewButton("OK")
	// ok.SetBorders(1, 1, 1, 1)
	// ok.SetPaddings(3, 3, 3, 3)
	// ok.SetMargins(3, 3, 3, 3)
	win.Add(cancel)
	cancel.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		app.Log().Info("Settings Cancelled")
		app.Gui().Root().SetModal(nil)
		app.mainPanel.Remove(win)
	})
	win.Add(ok)
	ok.Subscribe(gui.OnClick, func(e string, ev interface{}) {
		app.Log().Info("Settings Okayed")
		fmt.Printf("Debug: found ID: %s, chosen ID: %s\n", dDrop.Selected().Text(), tDrop.Selected().Text())
		sID, _ := strconv.Atoi(strings.Split(dDrop.Selected().Text(), ":")[0])
		err := openJoystick(sID, tDrop.Selected().Text())
		if err != nil {
			alertDialog(app.mainPanel, errorSev, err.Error())
		} else {
			app.settings.JoystickID = sID
			app.settings.JoystickType = tDrop.Selected().Text()
			if err = saveSettings(app.settings, appSettingsFile); err != nil {
				app.Log().Error(err.Error())
				alertDialog(app.mainPanel, errorSev, err.Error())
			}
		}
		app.Gui().Root().SetModal(nil)
		app.mainPanel.Remove(win)
	})

	root := app.mainPanel
	root.Add(win)
	win.SetPosition(root.Width()/2-settingsWidth/2, root.Height()/2-settingsHeight/2)
	app.Gui().SetModal(win)

	return win
}
