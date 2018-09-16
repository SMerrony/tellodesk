# Development Notes & Reminders

## Dependencies
* 3d0c/gmf for video handling. For Linux version to work gmf must be later than Sep 01 2018.
* simulatedsimian/joystick
* mattn/go-gtk
* SMerrony/tello >= v0.9.0

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
* Message Overlay updater - flightData.go:updateMessageCB()
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

Generated files:
* blueSky.gen.go
* icon.gen.go
  