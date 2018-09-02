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
* Video display updater
  * started in video.go:startVideo()
  * stopped in disconnectCB()

## Regularly-Run Funcs
* StatusBar updater - statusbar.go:updateStatusBarTCB()
  * Timer started in tdApp.go:setup() - 250ms
  * (No need to stop)
  