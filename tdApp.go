package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/SMerrony/tello"
	"github.com/g3n/engine/gui/assets/icon"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/application"
)

const (
	videoWidth, videoHeight = 960, 720
)

// tdApp holds GUI-related data, general data is currently globally defined in main()
type tdApp struct {
	*application.Application
	settingsLoaded                       bool
	settings                             settingsT
	menuBar                              *gui.Menu
	toolBar                              *toolbar
	mainPanel                            *gui.Panel
	tabBar                               *gui.TabBar
	statusBar                            *statusbar
	trackMenu, imagesMenu, flightSubMenu *gui.Menu     // just menus we need to access
	connectItem, disconnectItem          *gui.MenuItem // just the items we need to access
	tsmShowDrone, tsmShowPath            *gui.MenuItem
	recordVideoItem, stopRecordingItem   *gui.MenuItem
	importTrackItem                      *gui.MenuItem
	panel                                *gui.Panel
	feed                                 *gui.Image
	texture                              *texture.Texture2D
	picMu                                sync.RWMutex
	pic                                  *image.RGBA
	picChan                              chan *image.RGBA
	videoChan                            <-chan []byte
	videoRecMu                           sync.RWMutex
	videoRecording                       bool
	videoFile                            *os.File
	videoWriter                          *bufio.Writer
	flightDataMu                         sync.RWMutex
	flightData                           tello.FlightData
	trackChart                           *trackChartT
	trackTab                             *gui.Tab
	trackShowDrone, trackShowPath        bool
}

func (app *tdApp) setup() {

	//app.videoRecording.Store(false)

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
		if strings.Contains(err.Error(), "cannot find") || strings.Contains(err.Error(), "no such") {
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

	app.toolBar = app.buildToolbar()
	app.mainPanel.Add(app.toolBar)

	app.tabBar = gui.NewTabBar(videoWidth, videoHeight+20)
	app.mainPanel.Add(app.tabBar)

	//app.picChan = make(chan *image.RGBA, 1)

	app.buildFeed()
	feedTab := app.tabBar.AddTab("Feed")
	feedTab.SetPinned(true)
	feedTab.SetContent(app.feed)

	app.trackChart = buildTrackChart(videoWidth, videoHeight, defaultTrackScale, app.trackShowDrone, app.trackShowPath)
	app.trackTab = app.tabBar.AddTab("Tracker")
	app.trackTab.SetPinned(true)
	app.trackTab.SetContent(app.trackChart)

	planTab := app.tabBar.AddTab("Planner")
	planTab.SetPinned(true)

	app.tabBar.SetSelected(0)

	app.statusBar = buildStatusbar(app.mainPanel)
	app.mainPanel.Add(app.statusBar)
	//app.Subscribe("fdUpdate", app.updateStatusBar)
	app.Gui().TimerManager.SetInterval(250*time.Millisecond, true, app.updateStatusBarTCB)

	app.Gui().SetName(appName)

	app.Subscribe(application.OnQuit, app.exitNicely) // catch main window being closed
}

func (app *tdApp) buildMenu() {

	app.menuBar = gui.NewMenuBar()

	fileMenu := gui.NewMenu()
	settings := fileMenu.AddOption("Settings")
	settings.SetIcon(icon.Settings)
	settings.Subscribe(gui.OnClick, app.settingsCB)

	fileMenu.AddSeparator()

	ex := fileMenu.AddOption("Exit")
	ex.SetId("exit")
	ex.SetIcon(icon.Close)
	ex.Subscribe(gui.OnClick, app.exitNicely)

	app.menuBar.AddMenu("File ", fileMenu)

	droneMenu := gui.NewMenu()
	app.connectItem = droneMenu.AddOption("Connect")
	app.connectItem.SetIcon(icon.Sync)
	app.connectItem.Subscribe(gui.OnClick, app.connectCB)
	app.disconnectItem = droneMenu.AddOption("Disconnect")
	app.disconnectItem.SetIcon(icon.SyncDisabled)
	app.disconnectItem.SetEnabled(false).Subscribe(gui.OnClick, app.disconnectCB)

	app.flightSubMenu = gui.NewMenu()

	to := app.flightSubMenu.AddOption("Take-off")
	to.SetIcon(icon.FlightTakeoff)
	to.Subscribe(gui.OnClick, app.takeoffCB)
	tto := app.flightSubMenu.AddOption("Throw Take-off")
	tto.SetIcon(icon.ThumbUp)
	tto.Subscribe(gui.OnClick, app.throwTakeoffCB)
	lnd := app.flightSubMenu.AddOption("Land")
	lnd.SetIcon(icon.FlightLand)
	lnd.Subscribe(gui.OnClick, app.landCB)
	plnd := app.flightSubMenu.AddOption("Palm Land")
	plnd.SetIcon(icon.PanTool)
	plnd.Subscribe(gui.OnClick, app.palmLandCB)
	app.flightSubMenu.AddSeparator()
	sm := app.flightSubMenu.AddOption("Sports (Fast) Mode")
	sm.SetIcon(icon.DirectionsRun)
	sm.Subscribe(gui.OnClick, app.nyi)

	droneMenu.AddMenu("Flight", app.flightSubMenu)
	app.flightSubMenu.SetEnabled(false)

	app.menuBar.AddMenu(" Drone ", droneMenu)

	app.trackMenu = gui.NewMenu()
	ct := app.trackMenu.AddOption("Clear Track")
	ct.SetIcon(icon.Delete)
	ct.Subscribe(gui.OnClick, app.nyi)
	et := app.trackMenu.AddOption("Export Current Track as CSV")
	et.SetIcon(icon.Save)
	et.Subscribe(gui.OnClick, app.exportTrackCB)
	app.importTrackItem = app.trackMenu.AddOption("Import CSV Track")
	app.importTrackItem.SetIcon(icon.Input)
	app.importTrackItem.Subscribe(gui.OnClick, app.importTrackCB)
	st := app.trackMenu.AddOption("Save Track as PNG")
	st.SetIcon(icon.Image)
	st.Subscribe(gui.OnClick, app.exportTrackImageCB)

	trackSubMenu := gui.NewMenu()

	app.tsmShowDrone = trackSubMenu.AddOption("Show Drone Positions")
	app.tsmShowDrone.SetIcon(icon.CheckBox)
	app.trackShowDrone = true
	app.tsmShowDrone.Subscribe(gui.OnClick, app.trackShowDroneCB)
	app.tsmShowPath = trackSubMenu.AddOption("Show Track Path")
	app.tsmShowPath.SetIcon(icon.CheckBox)
	app.trackShowPath = true
	app.tsmShowPath.Subscribe(gui.OnClick, app.trackShowPathCB)

	app.trackMenu.AddMenu("Display", trackSubMenu)

	app.menuBar.AddMenu(" Track ", app.trackMenu)

	videoMenu := gui.NewMenu()
	app.recordVideoItem = videoMenu.AddOption("Record Video")
	app.recordVideoItem.SetIcon(icon.Videocam)
	app.recordVideoItem.Subscribe(gui.OnClick, app.recordVideoCB)
	app.stopRecordingItem = videoMenu.AddOption("Stop Recording")
	app.stopRecordingItem.SetIcon(icon.VideocamOff)
	app.stopRecordingItem.Subscribe(gui.OnClick, app.stopRecordingCB)
	app.stopRecordingItem.SetEnabled(false)
	app.menuBar.AddMenu(" Video ", videoMenu)

	app.imagesMenu = gui.NewMenu()
	tp := app.imagesMenu.AddOption("Take Photo")
	tp.SetIcon(icon.CameraAlt)
	tp.Subscribe(gui.OnClick, app.takePhotoCB)
	sp := app.imagesMenu.AddOption("Save Photo(s)")
	sp.SetIcon(icon.Save)
	sp.Subscribe(gui.OnClick, app.saveAllPhotosCB)
	app.menuBar.AddMenu(" Images ", app.imagesMenu)

	helpMenu := gui.NewMenu()
	oh := helpMenu.AddOption("Online Help")
	oh.SetIcon(icon.Help)
	oh.Subscribe(gui.OnClick, app.onlineHelpCB)
	helpMenu.AddSeparator()
	ab := helpMenu.AddOption("About")
	ab.SetIcon(icon.Info)
	ab.Subscribe(gui.OnClick, app.aboutCB)
	app.menuBar.AddMenu(" Help", helpMenu)

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
	app.pic = image.NewRGBA(image.Rect(0, 0, videoWidth, videoHeight))
}

func (app *tdApp) exitNicely(s string, i interface{}) {
	app.UnsubscribeID(application.OnQuit, nil) // prevent this being called again due to window app.Quit subscription
	app.Log().Info("Tidying-up and exiting")
	if drone.NumPics() > 0 {
		drone.SaveAllPics(fmt.Sprintf("tello_pic_%s", time.Now().Format("2006Jan2150405")))
	}
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

// func (app *tdApp) Render(a *app.App) {
// 	//app.statusBar.tm.ProcessTimers()
// }

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
