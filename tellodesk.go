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

	"github.com/SMerrony/tello"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//"net/http"
//_ "net/http/pprof"

const (
	appCopyright            = "©2018 S.Merrony"
	appDisclaimer           = "The author(s) is/are in no way\nconnected with Ryze®."
	appHelpURL              = "http://stephenmerrony.co.uk/blog/" // FIXME Help URL
	appIcon                 = "icon.png"                          // FIXME this should be a resource
	appName                 = "Tello® Desktop"
	appSettingsFile         = "tellodesktop.yaml"
	appVersion              = "0.1.0"
	bluesky                 = "sky960x720.png" // FIXME this should be a resource
	fdPeriodMs              = 100
	prefWidth, prefHeight   = videoWidth + 2, videoHeight + 72
	statusUpdatePeriodMs    = 250
	videoWidth, videoHeight = 960, 720
)

var appAuthors = []string{"Stephen Merrony"}

var (
	drone                                                 tello.Tello
	stickChan                                             chan<- tello.StickMessage
	jsStopChan, fdStopChan, vrStopChan, liveTrackStopChan chan bool
	fdChan                                                <-chan tello.FlightData
	videoChan                                             <-chan []byte
	stopFeedImageChan                                     chan bool
	feedWgt                                               *gtk.Image
	newFeedImageMu                                        sync.Mutex
	newFeedImage                                          bool
	feedImage                                             *image.RGBA
	videoRecMu                                            sync.RWMutex
	videoRecording                                        bool
	videoFile                                             *os.File
	videoWriter                                           *bufio.Writer
	win                                                   *gtk.Window
	menuBar                                               *menuBarT
	statusBar                                             *statusBarT

	flightDataMu sync.RWMutex
	flightData   tello.FlightData
	trackChart   *trackChartT

	settingsLoaded bool
	settings       settingsT
)

func main() {

	jsStopChan = make(chan bool) // not buffered
	fdStopChan = make(chan bool) // not buffered
	vrStopChan = make(chan bool) // not buffered
	liveTrackStopChan = make(chan bool)

	gtk.Init(nil)
	win = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.SetTitle(appName)
	win.SetIconFromFile(appIcon)
	getSettings()
	//win.SetDefaultSize(prefWidth, prefHeight)
	win.SetSizeRequest(prefWidth, prefHeight)
	win.SetResizable(false)
	win.Connect("destroy", func() {
		exitNicely()
	})

	vbox := gtk.NewVBox(false, 0)

	menuBar = buildMenu()
	vbox.PackStart(menuBar, false, false, 0)

	nb := gtk.NewNotebook()
	vbox.PackStart(nb, false, false, 1)

	feedWgt = buildFeedWgt()
	nb.AppendPage(feedWgt, gtk.NewLabel("Live Feed"))

	trackChart = buildTrackChart(videoWidth, videoHeight, defaultTrackScale,
		menuBar.trackShowDrone.GetActive(), menuBar.trackShowPath.GetActive())
	nb.AppendPage(trackChart, gtk.NewLabel("Tracker"))

	statusBar = buildStatusbar()
	vbox.PackEnd(statusBar, false, false, 0)
	glib.TimeoutAdd(statusUpdatePeriodMs, func() bool {
		statusBar.updateStatusBarTCB()
		return true
	})

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
		err = openJoystick(settings.JoystickID, settings.JoystickType)
		if err != nil {
			messageDialog(win, gtk.MESSAGE_ERROR, "Could not open configured joystick")
		}
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
	pb, _ := gdkpixbuf.NewPixbufFromFile(appIcon)
	about.SetLogo(pb)
	about.SetVersion(appVersion)
	about.SetAuthors(appAuthors)
	about.SetCopyright(appCopyright)
	about.SetComments(appDisclaimer)
	about.Run()
	about.Destroy()
}

func nyi() {
	messageDialog(win, gtk.MESSAGE_INFO, "Not Yet Implemented")
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
