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
	"log"
	"strconv"
	"time"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

type profileChartT struct {
	*gtk.Image
	track                                       *telloTrackT
	backingImage                                *image.RGBA
	pbd                                         gdkpixbuf.PixbufData
	pixBuf                                      *gdkpixbuf.Pixbuf
	width, height, xOrigin, yOrigin             int
	bgCol, axesCol, labelCol, lineCol, faintCol color.Color
	maxOffset                                   float32
	xScalePPS, yScalePPM                        float32 // scale factors expressed as Pixels Per Second/Metre
	trackDuration                               time.Duration
}

func buildProfileChart(w, h int, scale float32) (pc *profileChartT) {
	pc = new(profileChartT)
	pc.Image = gtk.NewImage()
	pc.width, pc.height = w, h
	pc.xOrigin = 20
	pc.yOrigin = h / 2
	pc.bgCol = color.White
	pc.axesCol = color.RGBA{0, 0, 0, 255}        // black
	pc.labelCol = color.RGBA{128, 128, 128, 255} // dark grey
	pc.lineCol = color.RGBA{255, 0, 0, 255}      // red
	pc.faintCol = color.RGBA{192, 192, 192, 64}  // light grey

	pc.maxOffset = scale
	pc.yScalePPM = float32(pc.yOrigin) / scale

	pc.trackDuration, _ = time.ParseDuration("1m")

	pc.backingImage = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	pc.pbd.Colorspace = gdkpixbuf.GDK_COLORSPACE_RGB
	pc.pbd.HasAlpha = true
	pc.pbd.BitsPerSample = 8
	pc.pbd.Width = w
	pc.pbd.Height = h
	pc.pbd.RowStride = pc.backingImage.Stride
	pc.pbd.Data = pc.backingImage.Pix
	pc.pixBuf = gdkpixbuf.NewPixbufFromData(pc.pbd)
	pc.track = newTrack()
	pc.drawEmptyChart()
	pc.SetFromPixbuf(pc.pixBuf)

	return pc
}

func (pc *profileChartT) clearChart() {
	draw.Draw(pc.backingImage, pc.backingImage.Bounds(), image.NewUniform(pc.bgCol), image.ZP, draw.Src)
	pc.pbd.Data = pc.backingImage.Pix
	pc.pixBuf = gdkpixbuf.NewPixbufFromData(pc.pbd)
	pc.SetFromPixbuf(pc.pixBuf)
}

func (pc *profileChartT) calcScales() {
	// vertical (height)
	pc.maxOffset = pc.track.deriveScale()
	pc.yScalePPM = float32(pc.yOrigin) / pc.maxOffset
	// horizontal (time)
	pc.trackDuration = pc.track.positions[len(pc.track.positions)-1].timeStamp.Sub(pc.track.positions[1].timeStamp) // FIXME
	pc.xScalePPS = float32(float64(pc.width-20) / pc.trackDuration.Seconds())
	log.Printf("Debug: profileChart xScalePPS is: %f, from %f seconds\n", pc.xScalePPS, pc.trackDuration.Seconds())
}

// xToOrd converts a horizontal (time in secs) value to its physical equivalent on an image
func (pc *profileChartT) xToOrd(x float32) (xOrd int) {
	xOrd = int(float32(pc.xOrigin) + x*pc.xScalePPS)
	return xOrd
}

// yToOrd converts a vertical (height) value to its physical equivalent on an image
func (pc *profileChartT) yToOrd(y float32) (yOrd int) {
	yOrd = int(float32(pc.yOrigin) - y*pc.yScalePPM)
	return yOrd
}

// drawEmptyChart draws the custom 'graph paper' for blank and populated charts.
func (pc *profileChartT) drawEmptyChart() {
	pc.clearChart()
	// blank vertical axis
	for y := 0; y < pc.height; y++ {
		pc.backingImage.Set(pc.xOrigin, y, pc.axesCol)
	}

	for s := 60; s <= int(pc.trackDuration.Seconds()); s += 60 {
		drawPhysLine(pc.backingImage, pc.xToOrd(float32(s)), 0, pc.xToOrd(float32(s)), pc.height, pc.faintCol)
	}

	// y-axis labels
	var yTickInterval float32 = 100.0
	switch {
	case pc.maxOffset < 10.1:
		yTickInterval = 1.0
	case pc.maxOffset < 101.0:
		yTickInterval = 10.0
	}
	for y := -pc.maxOffset; y <= pc.maxOffset; y += yTickInterval {
		pc.backingImage.Set(pc.xOrigin-1, pc.yOrigin+int(y*pc.yScalePPM), pc.axesCol)
		pc.backingImage.Set(pc.xOrigin+1, pc.yOrigin+int(y*pc.yScalePPM), pc.axesCol)
		drawPhysLabel(pc.backingImage, 5, pc.yToOrd(y)+6, strconv.Itoa(int(y)), pc.labelCol)
		drawPhysLine(pc.backingImage, pc.xOrigin, pc.yOrigin+int(y*pc.yScalePPM), pc.width, pc.yOrigin+int(y*pc.yScalePPM), pc.faintCol)
	}
	// blank horizontal axis
	for x := pc.xOrigin; x < pc.width; x++ {
		pc.backingImage.Set(x, pc.yOrigin, pc.axesCol)
	}
	pc.pbd.Data = pc.backingImage.Pix
	pc.pixBuf = gdkpixbuf.NewPixbufFromData(pc.pbd)
	pc.SetFromPixbuf(pc.pixBuf)
}

func (pc *profileChartT) drawTitles() {
	const dateFmt = "Jan 2 2006 15:04:05"
	if len(pc.track.positions) > 1 {
		drawPhysLabel(pc.backingImage, 40, pc.height-40, fmt.Sprintf("Flight Profile from %s to %s",
			pc.track.positions[1].timeStamp.Format(dateFmt),
			pc.track.positions[len(pc.track.positions)-1].timeStamp.Format("15:04:05")),
			pc.labelCol)
	}
}

func (pc *profileChartT) drawProfile() {
	if pc.track == nil {
		return
	}
	pc.calcScales()
	pc.drawEmptyChart()

	t0 := pc.track.positions[1].timeStamp
	var lastT time.Duration
	lastH := float32(pc.track.positions[1].heightDm) / 10
	for n, pos := range pc.track.positions {
		if n > 1 {
			t := pos.timeStamp.Sub(t0) // how long in (fractional) seconds
			h := float32(pos.heightDm) / 10
			// log.Printf("Debug: drawing from %d, %d to %d, %d\n", pc.xToOrd(float32(lastT.Seconds())), pc.yToOrd(lastH),
			// 	pc.xToOrd(float32(t.Seconds())), pc.yToOrd(h))
			drawPhysLine(pc.backingImage,
				pc.xToOrd(float32(lastT.Seconds())), pc.yToOrd(lastH),
				pc.xToOrd(float32(t.Seconds())), pc.yToOrd(h),
				pc.lineCol)
			lastT = t
			lastH = h
		}
	}

	pc.drawTitles()
	pc.pbd.Data = pc.backingImage.Pix
	pc.SetFromPixbuf(pc.pixBuf)

}
