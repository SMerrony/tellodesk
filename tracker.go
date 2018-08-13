package main

import (
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
