/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/SMerrony/tello"
	"github.com/mattn/go-gtk/gtk"
)

const timeStampFmt = "20060102150405.000"

// telloPosT defines an instantaneous position of the drone.
type telloPosT struct {
	timeStamp  time.Time
	heightDm   int16
	mvoX, mvoY float32
	imuYaw     int16
}

// telloTrackT contains a complete (or in-flight) track.
type telloTrackT struct {
	trackMu                sync.RWMutex
	maxX, maxY, minX, minY float32
	positions              []telloPosT
}

func newTrack() (tt *telloTrackT) {
	tt = new(telloTrackT)
	tt.positions = make([]telloPosT, 0, 1000)
	return tt
}

// toStrings converts a single position into an array of human-readable strings
// suitable for CSV export etc.
func (tp *telloPosT) toStrings() (strings []string) {
	strings = append(strings, tp.timeStamp.Format(timeStampFmt))
	strings = append(strings, fmt.Sprintf("%.3f", tp.mvoX))
	strings = append(strings, fmt.Sprintf("%.3f", tp.mvoY))
	strings = append(strings, fmt.Sprintf("%.1f", float64(tp.heightDm)/10))
	strings = append(strings, fmt.Sprintf("%d", tp.imuYaw))
	return strings
}

// toStruct does the inverse of toStrings, converting an array of strings into
// a single position struct.
func toStruct(strings []string) (tp telloPosT, err error) {
	tp.timeStamp, _ = time.Parse(timeStampFmt, strings[0])
	var f64 float64
	f64, _ = strconv.ParseFloat(strings[1], 32)
	tp.mvoX = float32(f64)
	f64, _ = strconv.ParseFloat(strings[2], 32)
	tp.mvoY = float32(f64)
	f64, _ = strconv.ParseFloat(strings[3], 32)
	tp.heightDm = int16(f64 * 10)
	i64, err := strconv.ParseInt(strings[4], 10, 16)
	tp.imuYaw = int16(i64)
	return tp, err
}

// addPositionIfChanged appends a new position report to the track if any of
// the mvoX, mvoY or Yaw have changed.  If only the timestamp has changed the
// position is not added.
func (tt *telloTrackT) addPositionIfChanged(fd tello.FlightData) {
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

// deriveScale returns the largest X or Y value rounded up to a whole number.
func (tt *telloTrackT) deriveScale() (scale float32) {
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

// simplify attempts to reduce the number of points in a track by eliminating consecutive postions that
// are within minDist metres of the previous position.
func (tt *telloTrackT) simplify(minDist float32) {
	if minDist < 0.1 || minDist > 2.0 {
		log.Printf("Track simplification ignored as minDist is out of range (0.1m ~ 2.0m): %f\n", minDist)
		return
	}
	if len(tt.positions) < 3 {
		log.Printf("Track too short to simplify (only %d points found).\n", len(tt.positions))
		return
	}
	min64 := float64(minDist)
	lastPos := tt.positions[0]
	thisIx := 1
	for thisIx < len(tt.positions) {
		thisPos := tt.positions[thisIx]
		xdiff := math.Abs(float64(lastPos.mvoX - thisPos.mvoX))
		ydiff := math.Abs(float64(lastPos.mvoY - thisPos.mvoY))
		zdiff := math.Abs(float64(lastPos.heightDm)-float64(thisPos.heightDm)) / 10.0
		if xdiff < min64 && ydiff < min64 && zdiff < min64 {
			tt.positions = append(tt.positions[:thisIx], tt.positions[thisIx+1:]...)
			//log.Printf("xdiff: %f, ydiff: %f, zdiff: %f ... skipping\n", xdiff, ydiff, zdiff)
		} else {
			lastPos = tt.positions[thisIx]
			thisIx++
			//log.Printf("xdiff: %f, ydiff: %f, zdiff: %f ... keeping\n", xdiff, ydiff, zdiff)
		}
	}
}

func simplifyCB() {
	posBefore := len(trackChart.track.positions)
	trackChart.track.simplify(0.3) // eliminates points within 30cm of each other
	posAfter := len(trackChart.track.positions)
	msg := fmt.Sprintf("Positions before : %d\n\nPositions after  : %d", posBefore, posAfter)
	messageDialog(win, gtk.MESSAGE_INFO, msg)
	trackChart.drawTrack()
}

// exportTrackCB exports the (global) current track as a CSV file.  The user is prompted for a filename.
func exportTrackCB() {
	var expPath string
	fs := gtk.NewFileChooserDialog(
		"File for Track Export",
		win,
		gtk.FILE_CHOOSER_ACTION_SAVE, "_Cancel", gtk.RESPONSE_CANCEL, "_Export", gtk.RESPONSE_ACCEPT)
	fs.SetCurrentFolder(settings.DataDir)
	ff := gtk.NewFileFilter()
	ff.AddPattern("*.csv")
	fs.SetFilter(ff)
	res := fs.Run()
	if res == gtk.RESPONSE_ACCEPT {
		expPath = fs.GetFilename()
		if expPath != "" {
			exp, err := os.Create(expPath)
			if err != nil {
				messageDialog(win, gtk.MESSAGE_INFO, "Could not create CSV file.")
			} else {
				defer exp.Close()
				w := csv.NewWriter(exp)
				trackChart.track.trackMu.RLock()
				for _, k := range trackChart.track.positions {
					w.Write(k.toStrings())
				}
				trackChart.track.trackMu.RUnlock()
				w.Flush()
			}
		}
	}
	fs.Destroy()
}

// exportTrackImageCB saves the currently-displayed track as a PNG image.  The user is prompted for a filename.
func exportTrackImageCB() {
	var expPath string
	fs := gtk.NewFileChooserDialog(
		"File for Track Image",
		win,
		gtk.FILE_CHOOSER_ACTION_SAVE, "_Cancel", gtk.RESPONSE_CANCEL, "_Export", gtk.RESPONSE_ACCEPT)
	fs.SetCurrentFolder(settings.DataDir)
	ff := gtk.NewFileFilter()
	ff.AddPattern("*.png")
	fs.SetFilter(ff)
	res := fs.Run()
	if res == gtk.RESPONSE_ACCEPT {
		expPath = fs.GetFilename()
		if expPath != "" {
			exp, err := os.Create(expPath)
			if err != nil {
				messageDialog(win, gtk.MESSAGE_INFO, "Could not create image file.")
			} else {
				defer exp.Close()
				if err := png.Encode(exp, trackChart.backingImage); err != nil {
					messageDialog(win, gtk.MESSAGE_INFO, "Could not write image file.")
				}
			}
		}
	}
	fs.Destroy()
}

// importTrackCB asks the user for the name of a CSV track and tries to import it via readTrack() as the current track.
func importTrackCB() {
	var impPath string
	fs := gtk.NewFileChooserDialog("Track to Import",
		win,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Import", gtk.RESPONSE_ACCEPT)
	fs.SetCurrentFolder(settings.DataDir)
	ff := gtk.NewFileFilter()
	ff.AddPattern("*.csv")
	fs.SetFilter(ff)
	res := fs.Run()
	if res == gtk.RESPONSE_ACCEPT {
		impPath = fs.GetFilename()
		if impPath != "" {
			imp, err := os.Open(impPath)
			if err != nil {
				messageDialog(win, gtk.MESSAGE_INFO, "Could not open track CSV file.")
			} else {
				defer imp.Close()
				r := csv.NewReader(bufio.NewReader(imp))
				tmpTrack := readTrack(r)
				trackChart.track = tmpTrack
				trackChart.drawTrack()
				profileChart.track = tmpTrack
				profileChart.drawProfile()
				notebook.SetCurrentPage(trackPage)
			}
		}
	}
	fs.Destroy()
}

// liveTracker is to be run at intervals (not as a goroutine)
func liveTrackerTCB() bool {
	trackChart.drawTrack()
	select {
	case <-liveTrackStopChan:
		return false
	default:
	}
	return true
}

// readTrack reads an open CSV file into a telloTrackT struct.
func readTrack(r *csv.Reader) (trk *telloTrackT) {
	trk = newTrack()
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			messageDialog(win, gtk.MESSAGE_INFO, "Could not read CSV track file.")
			return
		}
		tmpTrackPos, err := toStruct(line)
		if err != nil {
			messageDialog(win, gtk.MESSAGE_INFO, "Could not parse track CSV track file.")
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
	log.Printf("Imported %d track positions", len(trk.positions))
	log.Printf("Min X: %f, Max X:, %f", trk.minX, trk.maxX)
	log.Printf("Min Y: %f, Max Y:, %f", trk.minY, trk.maxY)
	log.Printf("Derived scale is %f", trk.deriveScale())
	return trk
}
