package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/texture"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// FixedLabel is a simple fixed-length, fixed-font textual label suitable
// for frequent updating.  The length is taken from the initial string.
type FixedLabel struct {
	gui.Panel
	contents      string             // value to be displayed
	width, height int                // label height in pixels
	length        int                // total # chars
	rgba          *image.RGBA        // backing image for the whole label
	tex           *texture.Texture2D // 2d textture for the whole label
	uniCol        *image.Uniform     // colour of drawn string
	rect          image.Rectangle    // precalculated label rect
}

// NewFixedLabel creates a fixed-size label which uses the built-in 7x13 font.
func NewFixedLabel(initial string, col color.Color) (l *FixedLabel) {
	l = new(FixedLabel)
	l.length = len(initial)
	l.height = 15
	l.width = 7 * l.length
	l.rect = image.Rectangle{image.Point{0, 0}, image.Point{l.width, l.height}}
	l.Panel.Initialize(float32(l.width), float32(l.height))
	l.rgba = image.NewRGBA(l.rect)
	l.tex = texture.NewTexture2DFromRGBA(l.rgba)
	l.tex.SetMagFilter(gls.NEAREST)
	l.tex.SetMinFilter(gls.NEAREST)
	l.Panel.Material().AddTexture(l.tex)
	l.uniCol = image.NewUniform(col)
	l.SetText(initial)
	return l
}

// SetText updates the text displayed on the FixedLabel.
func (l *FixedLabel) SetText(newString string) {
	// keep this as simple/fast as possible
	var text string
	if len(newString) > l.length {
		text = newString[:l.length]
	} else {
		text = newString
	}

	draw.Draw(l.rgba, l.rect, image.Transparent, image.ZP, draw.Src)
	d := &font.Drawer{
		Dst:  l.rgba,
		Src:  l.uniCol,
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(90), Y: fixed.Int26_6(750)},
	}
	d.DrawString(text)
	l.tex.SetFromRGBA(l.rgba)
}
