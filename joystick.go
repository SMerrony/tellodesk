/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/mattn/go-gtk/gtk"

	"github.com/SMerrony/tello"
	"github.com/simulatedsimian/joystick"
)

const (
	typeJoystick = iota
	typeGameController
	typeFlightController
)

const (
	axLeftX = iota
	axLeftY
	axRightX
	axRightY
	axLeftTrigger
	axRightTrigger
)

const (
	btnCross = iota
	btnCircle
	btnTriangle
	btnSquare
	btnA
	btnB
	btnBack
	btnL1
	btnL2
	btnL3
	btnR1
	btnR2
	btnR3
	btnStart
	btnX
	btnY
)

const deadZone = 2000

const jsUpdatePeriod = 40 * time.Millisecond // 40ms = 25Hz

// JoystickConfig holds a known joystick configuration
type JoystickConfig struct {
	Name    string
	JsType  int
	Axes    []int
	Buttons []uint
}

var (
	js                    joystick.Joystick
	jsID                  int
	jsConfig              JoystickConfig
	jsKnownWindowsConfigs = []JoystickConfig{
		JoystickConfig{
			Name:    "DualShock 3", // TODO - Untested
			JsType:  typeGameController,
			Axes:    []int{axLeftX: 0, axLeftY: 1, axRightX: 2, axRightY: 3},
			Buttons: []uint{btnCross: 1, btnCircle: 2, btnTriangle: 3, btnSquare: 0, btnL1: 4, btnL2: 6, btnR1: 5, btnR2: 7},
		},
		JoystickConfig{
			Name:    "DualShock 4",
			JsType:  typeGameController,
			Axes:    []int{axLeftX: 0, axLeftY: 1, axRightX: 2, axRightY: 3},
			Buttons: []uint{btnCross: 1, btnCircle: 2, btnTriangle: 3, btnSquare: 0, btnL1: 4, btnL2: 6, btnR1: 5, btnR2: 7},
		},
		JoystickConfig{
			Name:    "T-Flight Hotas X",
			JsType:  typeFlightController,
			Axes:    []int{axLeftX: 4, axLeftY: 2, axRightX: 0, axRightY: 1},
			Buttons: []uint{btnR1: 0, btnL1: 1, btnR3: 2, btnL3: 3, btnSquare: 4, btnCross: 5, btnCircle: 6, btnTriangle: 7, btnR2: 8, btnL2: 9},
		},
		JoystickConfig{
			Name:    "XBox 360", // TODO - Untested
			JsType:  typeGameController,
			Axes:    []int{axLeftX: 0, axLeftY: 1, axRightX: 3, axRightY: 4},
			Buttons: []uint{btnX: 2, btnY: 3, btnA: 0, btnB: 1, btnBack: 6, btnStart: 7},
		},
	}
	jsKnownLinuxConfigs = []JoystickConfig{
		JoystickConfig{
			Name:    "DualShock 4",
			JsType:  typeGameController,
			Axes:    []int{axLeftX: 0, axLeftY: 1, axRightX: 3, axRightY: 4},
			Buttons: []uint{btnCross: 0, btnCircle: 1, btnTriangle: 2, btnSquare: 3, btnL1: 4, btnL2: 6, btnR1: 5, btnR2: 7},
		},
		JoystickConfig{
			Name:    "T-Flight Hotas X", // Seeems to be the same on Linux and Windows
			JsType:  typeFlightController,
			Axes:    []int{axLeftX: 4, axLeftY: 2, axRightX: 0, axRightY: 1},
			Buttons: []uint{btnR1: 0, btnL1: 1, btnR3: 2, btnL3: 3, btnSquare: 4, btnCross: 5, btnCircle: 6, btnTriangle: 7, btnR2: 8, btnL2: 9},
		},
		JoystickConfig{
			Name:    "XBox 360", // TODO - Untested
			JsType:  typeGameController,
			Axes:    []int{axLeftX: 0, axLeftY: 1, axRightX: 3, axRightY: 4},
			Buttons: []uint{btnX: 2, btnY: 3, btnA: 0, btnB: 1, btnBack: 6, btnStart: 7},
		},
	}
)

// FoundJs holds one of the discovered joysticks
type FoundJs struct {
	ID   int
	Name string
}

func listJoysticks() (found []*FoundJs) {
	for jsid := 0; jsid < 10; jsid++ {
		js, err := joystick.Open(jsid)
		if err != nil {
			if jsid == 0 {
				fmt.Println("No joysticks detected")
				return nil
			}
		} else {
			fmt.Printf("Joystick ID: %d: Name: %s, Axes: %d, Buttons: %d\n", jsid, js.Name(), js.AxisCount(), js.ButtonCount())
			found = append(found, &FoundJs{jsid, fmt.Sprintf("%d: %s", jsid, js.Name())})
			js.Close()
		}
	}
	//fmt.Printf("Debug - listJoysticks returning: %v\n", found)
	return found
}

// KnownJs contains one of the known joystick types
type KnownJs struct {
	ID   int
	Name string
	Conf JoystickConfig
}

func listKnownJoystickTypes() (known []*KnownJs) {
	switch runtime.GOOS {
	case "windows":
		for jsid, config := range jsKnownWindowsConfigs {
			known = append(known, &KnownJs{jsid, config.Name, config})
		}
	case "linux":
		for jsid, config := range jsKnownLinuxConfigs {
			known = append(known, &KnownJs{jsid, config.Name, config})
		}
	}
	return known
}

func openJoystick(id int, chosenType string) (err error) {

	kt := listKnownJoystickTypes()
	for _, t := range kt {
		if t.Name == chosenType {
			jsConfig = t.Conf
			fmt.Printf("Debug: Joystick type set to: %s\n", jsConfig.Name)
			break
		}
	}

	js, err = joystick.Open(id)
	if err != nil {
		return errors.New("Could not open Joystick")
	}
	jsID = id

	return nil
}

func intAbs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

// readJoystick is run as a Goroutine
func readJoystick(test bool) {
	var (
		sm                 tello.StickMessage
		jsState, prevState joystick.State
		err                error
	)
	log.Println("Debug: Joystick listener starting")
	for {
		jsState, err = js.Read()

		if err != nil {
			log.Printf("Error reading joystick: %v\n", err)

		}

		if jsState.AxisData[jsConfig.Axes[axLeftX]] == 32768 {
			sm.Lx = 32767
		} else {
			sm.Lx = int16(jsState.AxisData[jsConfig.Axes[axLeftX]])
		}

		if jsState.AxisData[jsConfig.Axes[axLeftY]] == 32768 {
			sm.Ly = -32767
		} else {
			sm.Ly = -int16(jsState.AxisData[jsConfig.Axes[axLeftY]])
		}

		if jsState.AxisData[jsConfig.Axes[axRightX]] == 32768 {
			sm.Rx = 32767
		} else {
			sm.Rx = int16(jsState.AxisData[jsConfig.Axes[axRightX]])
		}

		if jsState.AxisData[jsConfig.Axes[axRightY]] == 32768 {
			sm.Ry = -32767
		} else {
			sm.Ry = -int16(jsState.AxisData[jsConfig.Axes[axRightY]])
		}

		if intAbs(sm.Lx) < deadZone {
			sm.Lx = 0
		}
		if intAbs(sm.Ly) < deadZone {
			sm.Ly = 0
		}
		if intAbs(sm.Rx) < deadZone {
			sm.Rx = 0
		}
		if intAbs(sm.Ry) < deadZone {
			sm.Ry = 0
		}

		if test {
			log.Printf("JS: Lx: %d, Ly: %d, Rx: %d=>%d, Ry: %d\n", sm.Lx, sm.Ly, jsState.AxisData[jsConfig.Axes[axRightX]], sm.Rx, sm.Ry)
		} else {
			stickChan <- sm
		}

		if jsState.Buttons&(1<<jsConfig.Buttons[btnL1]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnL1]) == 0 {
			if test {
				log.Println("L1 pressed")
			} else {
				drone.Bounce()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnL2]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnL2]) == 0 {
			if test {
				log.Println("L2 pressed")
			} else {
				drone.PalmLand()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnSquare]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnSquare]) == 0 {
			if test {
				log.Println("Square pressed")
			} else {
				drone.TakePicture()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnA]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnA]) == 0 {
			if test {
				log.Println("A pressed")
			} else {
				drone.TakePicture()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnTriangle]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnTriangle]) == 0 {
			if test {
				log.Println("Triangle pressed")
			} else {
				drone.TakeOff()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnY]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnY]) == 0 {
			if test {
				log.Println("Y pressed")
			} else {
				drone.TakeOff()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnCircle]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnCircle]) == 0 {
			if test {
				log.Println("Circle pressed")
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnCross]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnCross]) == 0 {
			if test {
				log.Println("Cross pressed")
			} else {
				drone.Land()
			}
		}
		if jsState.Buttons&(1<<jsConfig.Buttons[btnX]) != 0 && prevState.Buttons&(1<<jsConfig.Buttons[btnX]) == 0 {
			if test {
				log.Println("X pressed")
			} else {
				drone.Land()
			}
		}
		prevState = jsState

		time.Sleep(jsUpdatePeriod)
	}
}

func joystickHelpCB() {
	messageDialog(win, gtk.MESSAGE_INFO,
		`Joystick Controls

Right Stick   Forwards/backwards, left/right
Left Stick    Turn left/right, go up/down

▲ Triangle, Y (Yellow)   Take off
X  Cross, X (Blue)           Land
□ Square, A (Green)      Take Photo

L1 Button     Bounce
L2 Button     Palm Land
`)
}
