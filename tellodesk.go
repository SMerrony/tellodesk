/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/SMerrony/tello"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//"net/http"
//_ "net/http/pprof"

const (
	appCopyright  = "©2018 S.Merrony"
	appDisclaimer = "The author(s) is/are in no way\nconnected with Ryze®;\n" +
		"nor are they responsible for any\ndamage caused to, or by, any device\ncontrolled by this software."
	appHelpURL           = "https://github.com/SMerrony/tellodesk/wiki"
	appName              = "Tello® Desk"
	appSettingsFile      = "tellodesk.yaml"
	appVersion           = "v0.1.0"
	fdPeriodMs           = 100
	statusUpdatePeriodMs = 250
)

var appAuthors = []string{"Stephen Merrony"}

var (
	drone                                         tello.Tello
	stickChan                                     chan<- tello.StickMessage
	fdStopChan, vrStopChan, liveTrackStopChan     chan bool
	fdChan                                        <-chan tello.FlightData
	videoChan                                     <-chan []byte
	stopFeedImageChan                             chan bool
	videoWgt                                      *videoWgtT
	videoWidth, videoHeight                       = normalVideoWidth, normalVideoHeight
	win                                           *gtk.Window
	menuBar                                       *menuBarT
	notebook                                      *gtk.Notebook
	videoPage, statusPage, trackPage, profilePage int // IDs of the notebook pages for each tab
	statusBar                                     *statusBarT

	flightDataMu sync.RWMutex
	flightData   tello.FlightData
	statusTab    *liveStatusTabT
	liveTrack    *telloTrackT
	trackChart   *trackChartT
	profileChart *profileChartT

	settingsLoaded bool
	settings       settingsT

	blueSkyPixbuf, iconPixbuf *gdkpixbuf.Pixbuf
)

func main() {

	// preload the images from generated data
	blueSkyPixbuf = gdkpixbuf.NewPixbufFromData(blueSkyPNG)
	iconPixbuf = gdkpixbuf.NewPixbufFromData(iconPNG)

	fdStopChan = make(chan bool) // not buffered
	vrStopChan = make(chan bool) // not buffered
	liveTrackStopChan = make(chan bool)

	gtk.Init(nil)
	win = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.SetTitle(appName)
	win.SetIcon(iconPixbuf)
	getSettings()
	if settings.WideVideo {
		videoWidth, videoHeight = wideVideoWidth, wideVideoHeight
		blueSkyPixbuf = blueSkyPixbuf.ScaleSimple(videoWidth, videoHeight, gdkpixbuf.INTERP_BILINEAR)
	}
	win.SetResizable(false) // Gtk does the right thing and sets the size after laying out
	win.Connect("destroy", func() {
		exitNicely()
	})

	vbox := gtk.NewVBox(false, 0)

	menuBar = buildMenu()
	vbox.PackStart(menuBar, false, false, 0)

	notebook = gtk.NewNotebook()
	vbox.PackStart(notebook, false, false, 1)

	videoWgt = buildVideodWgt()
	videoPage = notebook.AppendPage(videoWgt, gtk.NewLabel("Live Feed"))

	statusTab = buildLiveStatusTab(videoWidth, videoHeight)
	statusPage = notebook.AppendPage(statusTab, gtk.NewLabel("Status"))

	liveTrack = newTrack()

	trackChart = buildTrackChart(liveTrack, videoWidth, videoHeight, defaultTrackScale,
		menuBar.trackShowDrone.GetActive(), menuBar.trackShowPath.GetActive())
	trackPage = notebook.AppendPage(trackChart, gtk.NewLabel("Tracker"))

	profileChart = buildProfileChart(videoWidth, videoHeight)
	profilePage = notebook.AppendPage(profileChart, gtk.NewLabel("Profile"))

	statusBar = buildStatusbar()
	vbox.PackEnd(statusBar, false, false, 0)
	glib.TimeoutAdd(statusUpdatePeriodMs, func() bool {
		statusBar.updateStatusBarTCB()
		return true
	})
	glib.TimeoutAdd(statusUpdatePeriodMs, updateFlightDataTCB)

	win.Add(vbox)
	win.ShowAll()
	gtk.Main()
}

func getSettings() {
	var err error
	settings, err = loadSettings(appSettingsFile)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find") || strings.Contains(err.Error(), "no such") {
			messageDialog(win, gtk.MESSAGE_INFO,
				"Could not open settings file\n\n"+appSettingsFile+"\n\n"+
					"This is normal on a first run,\nor until you have saved your settings")
		} else {
			messageDialog(win, gtk.MESSAGE_ERROR, err.Error())
		}
		settingsLoaded = false
		log.Printf("Error loading saved settings: %v", err)
	} else {
		log.Printf("Debug: loaded settings: chosen JS type is %s\n", settings.JoystickType)

		settingsLoaded = true
	}
}

func exitNicely() {
	log.Println("Tidying-up and exiting")
	if drone.NumPics() > 0 {
		saveAllPhotosCB()
	}
	gtk.MainQuit()
}

func aboutCB() {
	about := gtk.NewAboutDialog()
	about.SetProgramName(appName)
	about.SetIcon(iconPixbuf)
	about.SetLogo(iconPixbuf)
	about.SetVersion(appVersion)
	about.SetWebsite(appHelpURL)
	about.SetAuthors(appAuthors)
	about.SetCopyright(appCopyright)
	about.SetComments("Using Tello Package: " + tello.TelloPackageVersion + "\n\n" + appDisclaimer)
	about.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	about.Run()
	about.Destroy()
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
