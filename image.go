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
)

func (l *Layer) Bounds() image.Rectangle {
	x1, y1 := int(l.stl.Min.X), int(l.stl.Min.Y)
	x2, y2 := int(l.stl.Max.X+0.5), int(l.stl.Max.Y+0.5)
	return image.Rect(x1*drawfactor, y1*drawfactor, x2*drawfactor, y2*drawfactor)
}

func (l *Layer) Draw(dst draw.Image) {
	for _, s := range l.perimeters {
		drawLine(dst, perimeterColor, s)
	}
	for _, s := range l.infill {
		drawLine(dst, infillColor, s)
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
