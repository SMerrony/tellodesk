/**
 *Copyright (c) 2018 Stephen Merrony
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

type trackChartT struct {
	*gtk.Image
	track                                        *telloTrackT
	backingImage                                 *image.RGBA
	pbd                                          gdkpixbuf.PixbufData
	pixBuf                                       *gdkpixbuf.Pixbuf
	width, height, xOrigin, yOrigin              int
	bgCol, axesCol, labelCol, faintCol, droneCol color.Color
	maxOffset                                    float32
	scalePPM                                     float32 // scale factor expressed as Pixels Per Metre
	showDrone, showPath                          bool
}

const defaultTrackScale float32 = 10.0

func buildTrackChart(w, h int, scale float32, showDrone, showPath bool) (tc *trackChartT) {
	tc = new(trackChartT)
	tc.Image = gtk.NewImage()
	tc.width, tc.height = w, h
	tc.showDrone, tc.showPath = showDrone, showPath
	tc.xOrigin = w / 2
	tc.yOrigin = h / 2
	tc.bgCol = color.White
	tc.axesCol = color.RGBA{0, 0, 0, 255}        // black
	tc.labelCol = color.RGBA{128, 128, 128, 255} //dark grey
	tc.faintCol = color.RGBA{192, 192, 192, 64}  // light grey
	tc.droneCol = color.RGBA{255, 0, 0, 255}     // red
	tc.maxOffset = scale
	if w >= h { // scale to the shortest axis
		tc.scalePPM = float32(tc.yOrigin) / scale
	} else {
		tc.scalePPM = float32(tc.xOrigin) / scale
	}
	tc.backingImage = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	tc.pbd.Colorspace = gdkpixbuf.GDK_COLORSPACE_RGB
	tc.pbd.HasAlpha = true
	tc.pbd.BitsPerSample = 8
	tc.pbd.Width = w
	tc.pbd.Height = h
	tc.pbd.RowStride = tc.backingImage.Stride
	tc.pbd.Data = tc.backingImage.Pix
	tc.pixBuf = gdkpixbuf.NewPixbufFromData(tc.pbd)
	tc.track = newTrack()
	tc.drawEmptyChart()
	tc.SetFromPixbuf(tc.pixBuf)
	return tc
}

func (tc *trackChartT) calcScale() {
	tc.maxOffset = tc.track.deriveScale()
	if tc.width >= tc.height { // scale to the shortest axis
		tc.scalePPM = float32(tc.yOrigin) / tc.maxOffset
	} else {
		tc.scalePPM = float32(tc.xOrigin) / tc.maxOffset
	}
}

func (tc *trackChartT) clearChart() {
	draw.Draw(tc.backingImage, tc.backingImage.Bounds(), image.NewUniform(tc.bgCol), image.ZP, draw.Src)
	tc.pbd.Data = tc.backingImage.Pix
	tc.pixBuf = gdkpixbuf.NewPixbufFromData(tc.pbd)
	tc.SetFromPixbuf(tc.pixBuf)
}

// drawEmptyChart draws the custom 'graph paper' for blank and populated charts.
func (tc *trackChartT) drawEmptyChart() {
	tc.clearChart()
	// blank vertical axis
	drawPhysLine(tc.backingImage, tc.xOrigin, 0, tc.xOrigin, tc.height, tc.axesCol)

	// blank horizontal axis
	drawPhysLine(tc.backingImage, 0, tc.yOrigin, tc.width, tc.yOrigin, tc.axesCol)

	// x-axis labels
	var tickInterval float32 = 100.0
	switch {
	case tc.maxOffset < 10.1:
		tickInterval = 1.0
	case tc.maxOffset < 101.0:
		tickInterval = 10.0
	}
	for x := -tc.maxOffset; x <= tc.maxOffset; x += tickInterval {
		if x != 0 {
			drawPhysLine(tc.backingImage, tc.xToOrd(x), 0, tc.xToOrd(x), tc.height, tc.faintCol)
		}
		tc.backingImage.Set(tc.xOrigin+int(x*tc.scalePPM), tc.yOrigin-1, tc.axesCol)
		tc.backingImage.Set(tc.xOrigin+int(x*tc.scalePPM), tc.yOrigin+1, tc.axesCol)
		tc.drawLabel(x, 0, strconv.Itoa(int(x)))
	}
	// y-axis labels
	for y := -tc.maxOffset; y <= tc.maxOffset; y += tickInterval {
		if y != 0 {
			drawPhysLine(tc.backingImage, 0, tc.yToOrd(y), tc.width, tc.yToOrd(y), tc.faintCol)
		}
		tc.backingImage.Set(tc.xOrigin-1, tc.yOrigin+int(y*tc.scalePPM), tc.axesCol)
		tc.backingImage.Set(tc.xOrigin+1, tc.yOrigin+int(y*tc.scalePPM), tc.axesCol)
		tc.drawLabel(0, y, strconv.Itoa(int(y)))
	}
	tc.pbd.Data = tc.backingImage.Pix
	tc.pixBuf = gdkpixbuf.NewPixbufFromData(tc.pbd)
	tc.SetFromPixbuf(tc.pixBuf)
}

func (tc *trackChartT) drawLabel(x, y float32, lab string) {
	drawPhysLabel(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), lab, tc.labelCol)
}

func (tc *trackChartT) drawTitles() {
	const dateFmt = "Jan 2 2006 15:04:05"
	if len(tc.track.positions) > 1 {
		tc.drawLabel(-tc.maxOffset-0.5, -tc.maxOffset+0.5, fmt.Sprintf("Flight from %s to %s",
			tc.track.positions[1].timeStamp.Format(dateFmt),
			tc.track.positions[len(tc.track.positions)-1].timeStamp.Format("15:04:05")))
	}
}

// xToOrd converts a horizontal value to its physical equivalent on an image
func (tc *trackChartT) xToOrd(x float32) (xOrd int) {
	xOrd = int(float32(tc.xOrigin) + x*tc.scalePPM)
	return xOrd
}

// yToOrd converts a vertical value to its physical equivalent on an image
func (tc *trackChartT) yToOrd(y float32) (yOrd int) {
	yOrd = int(float32(tc.yOrigin) - y*tc.scalePPM)
	return yOrd
}

func (tc *trackChartT) drawPos(x, y float32, yaw int16) {
	switch {
	case yaw >= -45 && yaw <= 45: // N
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)+4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)+4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x), tc.yToOrd(y)+8, tc.droneCol)
	case yaw >= -135 && yaw < -45: // W
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)+4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)-4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+8, tc.yToOrd(y), tc.droneCol)
	case yaw > 45 && yaw < 135: // E
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)+4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)-4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-8, tc.yToOrd(y), tc.droneCol)
	default: // S
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)-4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)-4, tc.droneCol)
		drawPhysLine(tc.backingImage, tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x), tc.yToOrd(y)-8, tc.droneCol)
	}
}

func (tc *trackChartT) drawTrack() {
	if tc.track != nil {
		tc.calcScale()
	}
	tc.drawEmptyChart()
	var lastX, lastY float32
	for _, pos := range tc.track.positions {
		if tc.showDrone {
			tc.drawPos(pos.mvoX, pos.mvoY, pos.imuYaw)
		}
		if tc.showPath {
			tc.line(lastX, lastY, pos.mvoX, pos.mvoY, tc.droneCol)
			lastX = pos.mvoX
			lastY = pos.mvoY
		}
	}
	drawPhysLabel(tc.backingImage,
		tc.xToOrd(tc.track.positions[0].mvoX),
		tc.yToOrd(tc.track.positions[0].mvoY),
		"Start", tc.labelCol)
	last := len(tc.track.positions) - 1
	drawPhysLabel(tc.backingImage,
		tc.xToOrd(tc.track.positions[last].mvoX),
		tc.yToOrd(tc.track.positions[last].mvoY),
		"Finish", tc.labelCol)
	tc.drawTitles()
	tc.pbd.Data = tc.backingImage.Pix
	tc.SetFromPixbuf(tc.pixBuf)
}

// helper funcs...

func (tc *trackChartT) line(x0, y0, x1, y1 float32, col color.Color) {
	drawPhysLine(tc.backingImage, tc.xToOrd(x0), tc.yToOrd(y0), tc.xToOrd(x1), tc.yToOrd(y1), col)
}
