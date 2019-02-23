# Development Notes & Reminders
## Dependencies
* 3d0c/gmf for video handling. For Linux version to work gmf must be later than Sep 01 2018.
* simulatedsimian/joystick
* mattn/go-gtk
* SMerrony/tello >= v0.9.3

## Building on Ubuntu 18.04 (Bionic Beaver)
* Install all the -dev packages for libav* using any package manager
* `go get github.com/3d0c/gmf` - don't worry about the error
* cd into the github.com/3d0c/gmf directory
* `git checkout ec1401b491850f6cce7615f222072d0d473f1c80`
* redo the `go get` command above, there will be one warning you can safely ignore
  
## Func Naming Conventions
* Func names ending in ...CB are callbacks usually invoked from a menu or other GUI control
* Func names ending in ...TCB are timer callbacks to be regualarly run  via glib.TimeoutAdd() - they should return true for the timeout to be renewed.
* Important types are named ...T for clarity

## Goroutines
* Joystick reader 
  * started in droneCBs.go:connectCB(),
  * JS is closed in disconnectCB() which causes Goroutine to end
* FlightData listener 
  * started in droneCBs.go:connectCB(), 
  * stopped in disconnectCB()
* Video SPS/PPS Requestor 
  * started in video.go:startVideo()
  * stopped in disconnectCB()
* Video listener 
  * started in video.go:startVideo()

## Regularly-Run Funcs
* Video display updater
  * started in video.go:startVideo() - 30ms
  * stopped in disconnectCB()
* Flight Status updater - flightData.go:updateFlightDataTCB()
  * Started in main() - 250ms
  * (No need to stop)
* StatusBar updater - statusbar.go:updateStatusBarTCB()
  * Timer started in main - 250ms
  * (No need to stop)
* Live Tracker
  * Timer started in connectCB() - 500ms
  * Stopped in disconnectCB() via liveTrackStopChan

## Generated Files
Images are embedded using the go-gtk tool make_inline_pixbuf.  Command looks like:

`$GOBIN/make_inline_pixbuf.exe blueSkyPNG resources/sky960x720.png  > blueSky.gen.go`

 `make_inline_pixbuf.exe iconPNG resources/TD.png > icon.gen.go`

Generated files:
* blueSky.gen.go
* icon.gen.go
  