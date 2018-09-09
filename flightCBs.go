package main

func takeoffCB() {
	drone.TakeOff()
}
func throwTakeoffCB() {
	drone.ThrowTakeOff()
}
func landCB() {
	drone.Land()
}
func palmLandCB() {
	drone.PalmLand()
}

func toggleSportsModeCB() {
	drone.SetSportsMode(menuBar.sportsModeItem.GetActive())
}
