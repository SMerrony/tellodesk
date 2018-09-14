/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

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
