# Development Notes & Reminders
## Goroutines
* Joystick reader 
  * started in droneCBs.go:connectCB(), 
  * stopped in disconnectCB()
* FlightData listener 
  * started in droneCBs.go:connectCB(), 
  * stopped in disconnectCB()
* Video Restarter 
  * started in video.go:startVideo()
  * stopped in disconnectCB()
* Video listener 
  * started in video.go:startVideo()

## Regularly-Run Funcs
* Video display updater
  * started in video.go:startVideo() - 30ms
  * stopped in disconnectCB()
* ToolBar message updater
  * Timer started in main() - 250ms
  * (No need to stop)
* StatusBar updater - statusbar.go:updateStatusBarTCB()
  * Timer started in main - 250ms
  * (No need to stop)
* Live Tracker
  * Timer started in connectCB() - 500ms
  * Stopped in disconnectCB() via liveTrackStopChan
  