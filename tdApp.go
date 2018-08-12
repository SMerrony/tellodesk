package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/g3n/engine/gui/assets/icon"

	"github.com/g3n/g3nd/app"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/application"
)

const (
	videoWidth, videoHeight = 960, 720
)

type tdApp struct {
	*application.Application
	settingsLoaded                                                   bool
	settings                                                         settingsT
	menuBar                                                          *gui.Menu
	toolBar                                                          *toolbar
	mainPanel                                                        *gui.Panel
	statusBar                                                        *statusbar
	fileMenu, droneMenu, flightMenu, videoMenu, imagesMenu, helpMenu *gui.Menu
	connectItem, disconnectItem                                      *gui.MenuItem
	recordVideoItem, stopRecordingItem                               *gui.MenuItem
	panel                                                            *gui.Panel
	label                                                            *gui.Label
	feed                                                             *gui.Image
	texture                                                          *texture.Texture2D
	videoChan                                                        <-chan []byte
	videoRecording                                                   bool
	videoFile                                                        *os.File
	videoWriter                                                      *bufio.Writer
}

func (app *tdApp) setup() {
	// most stuff happens on the main panel
	app.mainPanel = gui.NewPanel(prefWidth, prefHeight)
	app.mainPanel.SetLayout(gui.NewVBoxLayout())
	app.Gui().Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		app.mainPanel.SetWidth(app.Gui().ContentWidth())
		app.mainPanel.SetHeight(app.Gui().ContentHeight())
		app.menuBar.SetWidth(app.Gui().ContentWidth())
	})
	app.Gui().Add(app.mainPanel)

	// load any saved settings now as they may affect the gui
	var err error
	app.settings, err = loadSettings(appSettingsFile)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find") {
			alertDialog(app.mainPanel, warningSev, "Could not open settings file\n\n"+appSettingsFile+"\n\n"+
				"This is normal on a first run,\nor until you have saved your settings")
		} else {
			alertDialog(app.mainPanel, warningSev, err.Error())
		}
		app.settingsLoaded = false
		app.Log().Info("Error loading saved settings: %v", err)
	} else {
		fmt.Printf("Debug: loaded settings: chosen JS type is %s\n", app.settings.JoystickType)
		err = openJoystick(app.settings.JoystickID, app.settings.JoystickType)
		if err != nil {
			alertDialog(app.mainPanel, errorSev, "Could not open configured joystick")
		}
		app.settingsLoaded = true
	}

	app.buildMenu()
	app.mainPanel.Add(app.menuBar)

	app.toolBar = buildToolbar(app.mainPanel)
	app.mainPanel.Add(app.toolBar)

	app.buildFeed()
	app.mainPanel.Add(app.feed)
	//app.feed.SetPosition(0, app.menuBar.Height())

	app.statusBar = buildStatusbar(app.mainPanel)
	app.mainPanel.Add(app.statusBar)

	app.Gui().SetName(appName)

	app.Subscribe(application.OnQuit, app.exitNicely) // catch main window being closed
}

func (app *tdApp) buildMenu() {
	app.menuBar = gui.NewMenuBar()
	app.fileMenu = gui.NewMenu()
	settings := app.fileMenu.AddOption("Settings")
	settings.SetIcon(icon.Settings)
	settings.Subscribe(gui.OnClick, app.settingsCB)

	app.fileMenu.AddSeparator()
	//app.fileMenu.AddOption("Exit").SetId("exit").Subscribe(gui.OnClick, func(s string, i interface{}) { app.Quit() })
	ex := app.fileMenu.AddOption("Exit")
	ex.SetId("exit")
	ex.SetIcon(icon.Close)
	ex.Subscribe(gui.OnClick, app.exitNicely)
	app.menuBar.AddMenu("File ", app.fileMenu)

	//app.menuBar.AddSeparator()

	app.droneMenu = gui.NewMenu()
	app.connectItem = app.droneMenu.AddOption("Connect")
	app.connectItem.SetIcon(icon.Sync)
	app.connectItem.Subscribe(gui.OnClick, app.connectCB)
	app.disconnectItem = app.droneMenu.AddOption("Disconnect")
	app.disconnectItem.SetIcon(icon.SyncDisabled)
	app.disconnectItem.SetEnabled(false).Subscribe(gui.OnClick, app.diconnectCB)
	app.menuBar.AddMenu(" Drone ", app.droneMenu)

	app.flightMenu = gui.NewMenu()
	to := app.flightMenu.AddOption("Take-off")
	to.SetIcon(icon.FlightTakeoff)
	to.Subscribe(gui.OnClick, app.takeoffCB)
	tto := app.flightMenu.AddOption("Throw Take-off")
	tto.SetIcon(icon.ThumbUp)
	tto.Subscribe(gui.OnClick, app.throwTakeoffCB)
	lnd := app.flightMenu.AddOption("Land")
	lnd.SetIcon(icon.FlightLand)
	lnd.Subscribe(gui.OnClick, app.landCB)
	plnd := app.flightMenu.AddOption("Palm Land")
	plnd.SetIcon(icon.PanTool)
	plnd.Subscribe(gui.OnClick, app.palmLandCB)
	app.flightMenu.AddSeparator()
	sm := app.flightMenu.AddOption("Sports (Fast) Mode")
	sm.SetIcon(icon.DirectionsRun)
	sm.Subscribe(gui.OnClick, app.nyi)
	app.menuBar.AddMenu(" Flight ", app.flightMenu)

	app.videoMenu = gui.NewMenu()
	app.recordVideoItem = app.videoMenu.AddOption("Record Video")
	app.recordVideoItem.SetIcon(icon.Videocam)
	app.recordVideoItem.Subscribe(gui.OnClick, app.recordVideoCB)
	app.stopRecordingItem = app.videoMenu.AddOption("Stop Recording")
	app.stopRecordingItem.SetIcon(icon.VideocamOff)
	app.stopRecordingItem.Subscribe(gui.OnClick, app.stopRecordingCB)
	app.stopRecordingItem.SetEnabled(false)
	app.menuBar.AddMenu(" Video ", app.videoMenu)

	app.imagesMenu = gui.NewMenu()
	tp := app.imagesMenu.AddOption("Take Photo")
	tp.SetIcon(icon.CameraAlt)
	tp.Subscribe(gui.OnClick, app.nyi)
	sp := app.imagesMenu.AddOption("Save Photo(s)")
	sp.SetIcon(icon.Save)
	sp.SetEnabled(false).Subscribe(gui.OnClick, app.nyi)
	app.menuBar.AddMenu(" Images ", app.imagesMenu)

	app.helpMenu = gui.NewMenu()
	oh := app.helpMenu.AddOption("Online Help")
	oh.SetIcon(icon.Help)
	oh.Subscribe(gui.OnClick, app.onlineHelpCB)
	app.helpMenu.AddSeparator()
	ab := app.helpMenu.AddOption("About")
	ab.SetIcon(icon.Info)
	ab.Subscribe(gui.OnClick, app.aboutCB)
	app.menuBar.AddMenu(" Help", app.helpMenu)

	app.menuBar.SetWidth(videoWidth)
}

func (app *tdApp) buildFeed() {
	const bluesky = "sky960x720.png"
	var err error
	app.texture, err = texture.NewTexture2DFromImage(bluesky)
	if err != nil {
		app.Log().Fatal("Could not load bluesky image - %v", err)
		app.Quit()
	}
	app.feed = gui.NewImageFromTex(app.texture)
}

func (app *tdApp) exitNicely(s string, i interface{}) {
	app.UnsubscribeID(application.OnQuit, nil) // prevent this being called again due to window app.Quit subscription
	app.Log().Info("Tidying-up and exiting")
	app.Quit()
}

func (app *tdApp) onlineHelpCB(s string, i interface{}) {
	openBrowser(appHelpURL)
}

func (app *tdApp) aboutCB(s string, i interface{}) {
	alertDialog(
		app.mainPanel,
		infoSev,
		fmt.Sprintf("Version: %s\n\nAuthor: %s\n\nCopyright: %s\n\nDisclaimer: %s", appVersion, appAuthor, appCopyright, appDisclaimer))
}

func (app *tdApp) nyi(s string, i interface{}) {
	alertDialog(app.mainPanel, infoSev, "Function not yet implemented")
}

func (app *tdApp) Render(a *app.App) {

}

// helper funcs

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
