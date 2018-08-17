package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/g3n/engine/gls"

	"github.com/g3n/engine/light"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	"github.com/g3n/engine/graphic"

	"github.com/SMerrony/tello"
	"github.com/g3n/engine/gui"
)

const timeStampFmt = "20060102150405.000"

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
	tt.positions = make([]telloPosT, 0, 1000)

	return tt
}

type trackChartT struct {
	*gui.Panel
	track *telloTrack
}

func (app *tdApp) buildTrackChart(w, h float32) (tc *trackChartT) {
	tc = new(trackChartT)
	tc.Panel = gui.NewPanel(w, h)
	//tc.Panel.SetRenderable(false)
	//tc.Panel.SetColor(math32.NewColor("green"))
	lab := gui.NewLabel("Hello")
	tc.Panel.Add(lab)
	app.SetPanel3D(tc.Panel)

	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewPhong(math32.NewColor("DarkBlue"))
	torusMesh := graphic.NewMesh(geom, mat)
	app.Scene().Add(torusMesh)

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	app.Scene().Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	app.CameraPersp().SetPosition(0, 0, 3)

	// // **** Just testing ****
	// // Adds ambient light to the test scene
	// ambLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)
	// app.Scene().Add(ambLight)

	// // Sets perspective camera position
	// width, height := app.Window().Size()
	// aspect := float32(width) / float32(height)
	// app.CameraPersp().SetPosition(0, 0, 5)
	// app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})
	// app.CameraPersp().SetAspect(aspect)

	// // Sets orthographic camera
	// app.CameraOrtho().SetPosition(0, 0, 3)
	// app.CameraOrtho().LookAt(&math32.Vector3{0, 0, 0})
	// app.CameraOrtho().SetZoom(1.0)

	// // Default camera is perspective
	// app.SetCamera(app.CameraPersp())
	// // Adds camera to scene (important for audio demos)
	// app.Scene().Add(app.Camera().GetCamera())

	// // **** End of Test Stuff ****

	tc.track = newTrack()
	return tc
}

func (tc *trackChartT) Render(*gls.GLS) {

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
	var pos telloPosT

	pos.heightDm = fd.Height
	pos.mvoX = fd.MVO.PositionX
	pos.mvoY = fd.MVO.PositionY
	pos.imuYaw = fd.IMU.Yaw

	if len(tt.positions) == 0 {
		tt.trackMu.Lock()
		tt.positions = append(tt.positions, pos)
		tt.trackMu.Unlock()
	} else {
		lastPos := tt.positions[len(tt.positions)-1]
		if lastPos.heightDm != pos.heightDm || lastPos.mvoX != pos.mvoX || lastPos.mvoY != pos.mvoY || lastPos.imuYaw != pos.imuYaw {
			pos.timeStamp = time.Now()
			tt.trackMu.Lock()
			tt.positions = append(tt.positions, pos)
			tt.trackMu.Unlock()
		}
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
				w.Flush()
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
	fs, _ := NewFileSelect(app.mainPanel, cwd, "Choose CSV Path for Import", ".csv")
	fs.Subscribe("OnOK", func(n string, ev interface{}) {
		impPath = fs.Selected()
		if impPath != "" {
			imp, err := os.Open(impPath)
			if err != nil {
				alertDialog(app.mainPanel, warningSev, "Could not open CSV file")
			} else {
				defer imp.Close()
				r := csv.NewReader(bufio.NewReader(imp))
				app.trackChart.track = newTrack()
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
					app.trackChart.track.positions = append(app.trackChart.track.positions, tmpTrackPos)
				}
				app.Log().Info("Imported %d track positions", len(app.trackChart.track.positions))
			}
		}
		fs.Close()
	})
	fs.Subscribe("OnCancel", func(n string, ev interface{}) {
		fs.Close()
	})
}
