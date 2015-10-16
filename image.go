package slice

import (
	"image"
	"image/color"
	"image/draw"

	"sigint.ca/graphics/primitive"
)

const drawfactor = 20

var (
	perimeterColor = color.Black
	infillColor    = color.RGBA{R: 0xFF, G: 0, B: 0, A: 0xFF}
	debugColor     = color.RGBA{R: 0, G: 0xFF, B: 0, A: 0xFF}
)

// Bounds returns an image.Rectangle for use with Draw.
func (l *Layer) Bounds() image.Rectangle {
	x1, y1 := int(l.stl.Min.X*drawfactor), int(l.stl.Min.Y*drawfactor)
	x2, y2 := int(l.stl.Max.X*drawfactor+0.5), int(l.stl.Max.Y*drawfactor+0.5)
	return image.Rect(x1, y1, x2, y2)
}

// Draw draws an image representation of the layer onto dst. Use Bounds to find
// the minimum size that dst should be to contain the entire image.
func (l *Layer) Draw(dst draw.Image) {
	for _, s := range l.perimeters {
		drawLine(dst, perimeterColor, s)
	}
	for _, s := range l.infill {
		drawLine(dst, infillColor, s)
	}
	for _, s := range l.debug {
		drawLine(dst, debugColor, s)
	}
}

func drawLine(dst draw.Image, c color.Color, seg *segment) {
	p1 := image.Pt(round(seg.from.X*drawfactor), round(seg.from.Y*drawfactor))
	p2 := image.Pt(round(seg.to.X*drawfactor), round(seg.to.Y*drawfactor))
	primitive.Line(dst, c, p1, p2)
}

func round(v float64) int {
	if v > 0.0 {
		return int(v + 0.5)
	} else if v < 0.0 {
		return int(v - 0.5)
	} else {
		return 0
	}
}
