package main

func (app *tdApp) takeoffCB(s string, i interface{}) {
	drone.TakeOff()
}
func (app *tdApp) throwTakeoffCB(s string, i interface{}) {
	drone.ThrowTakeOff()
}
func (app *tdApp) landCB(s string, i interface{}) {
	drone.Land()
}
func (app *tdApp) palmLandCB(s string, i interface{}) {
	drone.PalmLand()
}
