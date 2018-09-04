package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/texture"
)

type trackChartT struct {
	gui.Panel
	track                              *telloTrack
	tex                                *texture.Texture2D
	backingImage                       *image.RGBA
	width, height, xOrigin, yOrigin    int
	bgCol, axesCol, labelCol, droneCol color.Color
	maxOffset                          float32
	scalePPM                           float32 // scale factor expressed as Pixels Per Metre
	showDrone, showPath                bool
}

const defaultTrackScale float32 = 10.0

func buildTrackChart(w, h int, scale float32, showDrone, showPath bool) (tc *trackChartT) {
	tc = new(trackChartT)
	tc.width = w
	tc.height = h
	tc.showDrone, tc.showPath = showDrone, showPath
	tc.Panel.Initialize(float32(tc.width), float32(tc.height))
	tc.xOrigin = w / 2
	tc.yOrigin = h / 2
	tc.bgCol = color.White
	tc.axesCol = color.RGBA{128, 128, 128, 255} // color.Black
	tc.labelCol = color.RGBA{128, 128, 128, 255}
	tc.droneCol = color.RGBA{255, 0, 0, 255}
	tc.maxOffset = scale
	if w >= h { // scale to the shortest axis
		tc.scalePPM = float32(tc.yOrigin) / scale
	} else {
		tc.scalePPM = float32(tc.xOrigin) / scale
	}
	tc.backingImage = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	tc.tex = texture.NewTexture2DFromRGBA(tc.backingImage)
	tc.tex.SetMagFilter(gls.NEAREST)
	tc.tex.SetMinFilter(gls.NEAREST)
	tc.Panel.Material().AddTexture(tc.tex)
	tc.track = newTrack()
	tc.drawEmptyChart()
	return tc
}

func (tc *trackChartT) clearChart() {
	draw.Draw(tc.backingImage, tc.backingImage.Bounds(), image.NewUniform(tc.bgCol), image.ZP, draw.Src)
	tc.tex.SetFromRGBA(tc.backingImage)
}

func (tc *trackChartT) setShowDrone(show bool) {
	tc.showDrone = show
}

func (tc *trackChartT) setShowPath(show bool) {
	tc.showPath = show
}

func (tc *trackChartT) drawEmptyChart() {
	tc.tex.Dispose()
	tc.clearChart()
	// blank vertical axis
	for y := 0; y < tc.height; y++ {
		tc.backingImage.Set(tc.xOrigin, y, tc.axesCol)
	}
	// blank horizontal axis
	for x := 0; x < tc.width; x++ {
		tc.backingImage.Set(x, tc.yOrigin, tc.axesCol)
	}
	// x-axis labels
	var tickInterval float32 = 100.0
	switch {
	case tc.maxOffset < 10.1:
		tickInterval = 1.0
	case tc.maxOffset < 101.0:
		tickInterval = 10.0
	}
	for x := -tc.maxOffset; x <= tc.maxOffset; x += tickInterval {
		tc.backingImage.Set(tc.xOrigin+int(x*tc.scalePPM), tc.yOrigin-1, tc.axesCol)
		tc.backingImage.Set(tc.xOrigin+int(x*tc.scalePPM), tc.yOrigin+1, tc.axesCol)
		tc.drawLabel(x, 0, strconv.Itoa(int(x)))
	}
	// y-axis labels
	for y := -tc.maxOffset; y <= tc.maxOffset; y += tickInterval {
		tc.backingImage.Set(tc.xOrigin-1, tc.yOrigin+int(y*tc.scalePPM), tc.axesCol)
		tc.backingImage.Set(tc.xOrigin+1, tc.yOrigin+int(y*tc.scalePPM), tc.axesCol)
		tc.drawLabel(0, y, strconv.Itoa(int(y)))
		//fmt.Printf("Y label drawn at: %f\n", y)
	}
	tc.tex.SetFromRGBA(tc.backingImage)
}

func (tc *trackChartT) drawLabel(x, y float32, lab string) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(tc.xToOrd(x) * 64),
		Y: fixed.Int26_6(tc.yToOrd(y) * 64)}
	d := &font.Drawer{
		Dst:  tc.backingImage,
		Src:  image.NewUniform(tc.labelCol),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(lab)
}

func (tc *trackChartT) drawTitles() {
	const dateFmt = "Jan 2 2006 15:04:05"
	if len(tc.track.positions) > 1 {
		tc.drawLabel(-tc.maxOffset-0.5, -tc.maxOffset+0.5, fmt.Sprintf("Flight from %s to %s",
			tc.track.positions[1].timeStamp.Format(dateFmt),
			tc.track.positions[len(tc.track.positions)-1].timeStamp.Format("15:04:05")))
	}
}

func (tc *trackChartT) xToOrd(x float32) (xOrd int) {
	xOrd = int(float32(tc.xOrigin) + x*tc.scalePPM)
	return xOrd
}

func (tc *trackChartT) yToOrd(y float32) (yOrd int) {
	yOrd = int(float32(tc.yOrigin) - y*tc.scalePPM)
	return yOrd
}

func (tc *trackChartT) drawPos(x, y float32, yaw int16) {
	switch {
	case yaw >= -45 && yaw <= 45: // N
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)+4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)+4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x), tc.yToOrd(y)+8, tc.droneCol)
	case yaw >= -135 && yaw < -45: // W
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)+4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)-4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-8, tc.yToOrd(y), tc.droneCol)
	case yaw > 45 && yaw < 135: // E
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)+4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)-4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+8, tc.yToOrd(y), tc.droneCol)
	default: // S
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)-4, tc.yToOrd(y)-4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x)+4, tc.yToOrd(y)-4, tc.droneCol)
		tc.physLine(tc.xToOrd(x), tc.yToOrd(y), tc.xToOrd(x), tc.yToOrd(y)-8, tc.droneCol)
	}
}

func (tc *trackChartT) drawTrack() {
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
	tc.drawTitles()
}

// helper funcs...

func (tc *trackChartT) line(x0, y0, x1, y1 float32, col color.Color) {
	tc.physLine(tc.xToOrd(x0), tc.yToOrd(y0), tc.xToOrd(x1), tc.yToOrd(y1), col)
}

func (tc *trackChartT) physLine(x0, y0, x1, y1 int, col color.Color) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy
	for {
		tc.backingImage.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
