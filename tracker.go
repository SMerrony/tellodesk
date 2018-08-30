package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/SMerrony/tello"
)

const timeStampFmt = "20060102150405.000"

type telloPosT struct {
	timeStamp  time.Time
	heightDm   int16
	mvoX, mvoY float32
	imuYaw     int16
}

type telloTrack struct {
	trackMu                sync.RWMutex
	maxX, maxY, minX, minY float32
	positions              []telloPosT
}

func newTrack() (tt *telloTrack) {
	tt = new(telloTrack)
	tt.positions = make([]telloPosT, 0, 1000)

	return tt
}

func (tp *telloPosT) toStrings() (strings []string) {
	strings = append(strings, tp.timeStamp.Format(timeStampFmt))
	strings = append(strings, fmt.Sprintf("%.3f", tp.mvoX))
	strings = append(strings, fmt.Sprintf("%.3f", tp.mvoY))
	strings = append(strings, fmt.Sprintf("%.1f", float64(tp.heightDm)/10))
	strings = append(strings, fmt.Sprintf("%d", tp.imuYaw))
	return strings
}

func toStruct(strings []string) (tp telloPosT, err error) {
	tp.timeStamp, err = time.Parse(timeStampFmt, strings[0])
	var f64 float64
	f64, err = strconv.ParseFloat(strings[1], 32)
	tp.mvoX = float32(f64)
	f64, err = strconv.ParseFloat(strings[2], 32)
	tp.mvoY = float32(f64)
	f64, err = strconv.ParseFloat(strings[3], 32)
	tp.heightDm = int16(f64 * 10)
	i64, err := strconv.ParseInt(strings[4], 10, 16)
	tp.imuYaw = int16(i64)
	return tp, err
}

func (tt *telloTrack) addPositionIfChanged(fd tello.FlightData) {
	var newPos telloPosT

	newPos.heightDm = fd.Height
	newPos.mvoX = fd.MVO.PositionX
	newPos.mvoY = fd.MVO.PositionY
	newPos.imuYaw = fd.IMU.Yaw

	if len(tt.positions) == 0 {
		tt.trackMu.Lock()
		tt.positions = append(tt.positions, newPos)
		tt.trackMu.Unlock()
	} else {
		lastPos := tt.positions[len(tt.positions)-1]
		if lastPos.heightDm == newPos.heightDm && lastPos.mvoX == newPos.mvoX && lastPos.mvoY == newPos.mvoY && lastPos.imuYaw == newPos.imuYaw {
			// nothing has changed - just return
			return
		}
		newPos.timeStamp = time.Now()
		tt.trackMu.Lock()
		tt.positions = append(tt.positions, newPos)
		tt.trackMu.Unlock()
	}

	if newPos.mvoX < tt.minX {
		tt.minX = newPos.mvoX
	}
	if newPos.mvoX > tt.maxX {
		tt.maxX = newPos.mvoX
	}
	if newPos.mvoY < tt.minY {
		tt.minY = newPos.mvoY
	}
	if newPos.mvoY > tt.maxY {
		tt.maxY = newPos.mvoY
	}
}

func (app *tdApp) exportTrackCB(s string, ev interface{}) {
	var expPath string
	cwd, _ := os.Getwd()
	fs, _ := NewFileSelect(app.mainPanel, cwd, "Choose File for Track Export", "*.csv")
	fs.Subscribe("OnOK", func(n string, ev interface{}) {
		expPath = fs.Selected()
		if expPath != "" {
			exp, err := os.Create(expPath)
			if err != nil {
				alertDialog(app.mainPanel, warningSev, "Could not create CSV file")
			} else {
				defer exp.Close()
				w := csv.NewWriter(exp)
				app.trackChart.track.trackMu.RLock()
				for _, k := range app.trackChart.track.positions {
					w.Write(k.toStrings())
				}
				app.trackChart.track.trackMu.RUnlock()
				w.Flush()
			}
		}
		fs.Close()
	})
	fs.Subscribe("OnCancel", func(n string, ev interface{}) {
		fs.Close()
	})
}

func (app *tdApp) exportTrackImageCB(s string, ev interface{}) {
	var expPath string
	cwd, _ := os.Getwd()
	fs, _ := NewFileSelect(app.mainPanel, cwd, "Choose File for Track Image", "*.png")
	fs.Subscribe("OnOK", func(n string, ev interface{}) {
		expPath = fs.Selected()
		if expPath != "" {
			exp, err := os.Create(expPath)
			if err != nil {
				alertDialog(app.mainPanel, warningSev, "Could not create image file")
			} else {
				defer exp.Close()
				if err := png.Encode(exp, app.trackChart.backingImage); err != nil {
					alertDialog(app.mainPanel, errorSev, "Could not export track image")
				}
			}
		}
		fs.Close()
	})
	fs.Subscribe("OnCancel", func(n string, ev interface{}) {
		fs.Close()
	})
}

func (app *tdApp) importTrackCB(s string, ev interface{}) {
	var impPath string
	cwd, _ := os.Getwd()
	fs, _ := NewFileSelect(app.mainPanel, cwd, "Choose CSV Path for Import", "*.csv")
	fs.Subscribe("OnOK", func(n string, ev interface{}) {
		impPath = fs.Selected()
		if impPath != "" {
			imp, err := os.Open(impPath)
			if err != nil {
				alertDialog(app.mainPanel, warningSev, "Could not open CSV file")
			} else {
				defer imp.Close()
				r := csv.NewReader(bufio.NewReader(imp))
				tmpTrack := app.readTrack(r)
				app.trackChart = buildTrackChart(videoWidth, videoHeight, tmpTrack.deriveScale(), app.trackShowDrone, app.trackShowPath)
				app.trackChart.track = tmpTrack
				app.trackChart.drawTrack()
				app.trackTab.SetContent(app.trackChart)
			}
		}
		fs.Close()
	})
	fs.Subscribe("OnCancel", func(n string, ev interface{}) {
		fs.Close()
	})
}

func (app *tdApp) readTrack(r *csv.Reader) (trk *telloTrack) {
	trk = newTrack()
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			alertDialog(app.mainPanel, errorSev, "Could not read CSV file")
			return
		}
		tmpTrackPos, err := toStruct(line)
		if err != nil {
			alertDialog(app.mainPanel, errorSev, "Could not parse CSV file")
			return
		}
		trk.positions = append(trk.positions, tmpTrackPos)

		if tmpTrackPos.mvoX < trk.minX {
			trk.minX = tmpTrackPos.mvoX
		}
		if tmpTrackPos.mvoX > trk.maxX {
			trk.maxX = tmpTrackPos.mvoX
		}
		if tmpTrackPos.mvoY < trk.minY {
			trk.minY = tmpTrackPos.mvoY
		}
		if tmpTrackPos.mvoY > trk.maxY {
			trk.maxY = tmpTrackPos.mvoY
		}

	}
	app.Log().Info("Imported %d track positions", len(trk.positions))
	app.Log().Info("Min X: %f, Max X:, %f", trk.minX, trk.maxX)
	app.Log().Info("Min Y: %f, Max Y:, %f", trk.minY, trk.maxY)
	app.Log().Info("Derived scale is %f", trk.deriveScale())
	return trk
}

func (tt *telloTrack) deriveScale() (scale float32) {
	scale = 1.0 // minimum scale value
	if tt.maxX > scale {
		scale = tt.maxX
	}
	if -tt.minX > scale {
		scale = -tt.minX
	}
	if tt.maxY > scale {
		scale = tt.maxY
	}
	if -tt.minY > scale {
		scale = -tt.minY
	}
	scale = float32(math.Ceil(float64(scale)))
	return scale
}
