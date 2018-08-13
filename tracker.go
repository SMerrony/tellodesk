package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/SMerrony/tello"
)

type telloPosT struct {
	timeStamp  time.Time
	heightDm   int16
	mvoX, mvoY float32
	imuYaw     int16
}

type telloTrack struct {
	trackMu            sync.RWMutex
	startTime, endTime time.Time
	positions          []telloPosT
}

func newTrack() (tt *telloTrack) {
	tt = new(telloTrack)
	tt.positions = make([]telloPosT, 1000)

	return tt
}

func (tp *telloPosT) toStrings() (strings []string) {
	strings = append(strings, tp.timeStamp.String())
	strings = append(strings, fmt.Sprintf("%.1f", float64(tp.heightDm)/10))
	strings = append(strings, fmt.Sprintf("%.3f", tp.mvoX))
	strings = append(strings, fmt.Sprintf("%.3f", tp.mvoY))
	strings = append(strings, fmt.Sprintf("%d", tp.imuYaw))
	return strings
}

func (tt *telloTrack) addPositionIfChanged(fd tello.FlightData) {
	var pos telloPosT

	pos.heightDm = fd.Height
	pos.mvoX = fd.MVO.PositionX
	pos.mvoY = fd.MVO.PositionY
	pos.imuYaw = fd.IMU.Yaw

	lastPos := tt.positions[len(tt.positions)-1]
	if lastPos.heightDm != pos.heightDm || lastPos.mvoX != pos.mvoX || lastPos.mvoY != pos.mvoY || lastPos.imuYaw != pos.imuYaw {
		pos.timeStamp = time.Now()
		tt.trackMu.Lock()
		tt.positions = append(tt.positions, pos)
		tt.trackMu.Unlock()
	}
}

func (app *tdApp) exportTrackCB(s string, ev interface{}) {
	var expPath string
	cwd, _ := os.Getwd()
	fs, _ := NewFileSelect(app.mainPanel, cwd, "Choose File for Path Export", ".csv")
	fs.Subscribe("OnOK", func(n string, ev interface{}) {
		expPath = fs.Selected()
		if expPath != "" {
			exp, err := os.Create(expPath)
			if err != nil {
				alertDialog(app.mainPanel, warningSev, "Could not create CSV file")
			} else {
				defer exp.Close()
				w := csv.NewWriter(exp)
				currentTrack.trackMu.RLock()
				for _, k := range currentTrack.positions {
					w.Write(k.toStrings())
				}
				currentTrack.trackMu.RUnlock()
			}
		}
		fs.Close()
	})
	fs.Subscribe("OnCancel", func(n string, ev interface{}) {
		fs.Close()
	})
}
