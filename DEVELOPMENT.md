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
* Video display image updater - video.go:updateFeedTCB()
  * Timer started in droneCBs.go:connectCB() - 33ms
  * Timer stopped in droneCBs.go:disconnectCB()
* StatusBar updater - statusbar.go:updateStatusBarTCB()
  * Timer started in tdApp.go:setup() - 250ms
  * (No need to stop)
  