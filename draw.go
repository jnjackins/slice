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
	x1, y1 := int(l.stl.Min.X*drawfactor-0.5), int(l.stl.Min.Y*drawfactor-0.5)
	x2, y2 := int(l.stl.Max.X*drawfactor+0.5), int(l.stl.Max.Y*drawfactor+0.5)
	return image.Rect(x1, y1, x2, y2)
}

// Draw draws an image representation of the layer onto dst. Use Bounds to find
// the minimum size that dst should be to contain the entire image.
func (l *Layer) Draw(dst draw.Image) {
	for _, solid := range l.solids {
		for _, s := range solid.perimeters {
			drawLine(dst, s.from, s.to, perimeterColor)
		}
		for _, s := range solid.infill {
			drawLine(dst, s.from, s.to, infillColor)
		}
		for _, s := range solid.debug {
			drawLine(dst, s.from, s.to, debugColor)
		}
	}
}

func drawLine(dst draw.Image, p1, p2 Vertex2, c color.Color) {
	px1 := image.Pt(round(p1.X*drawfactor), round(p1.Y*drawfactor))
	px2 := image.Pt(round(p2.X*drawfactor), round(p2.Y*drawfactor))
	primitive.Line(dst, c, px1, px2)
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
