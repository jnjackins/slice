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

func (l *Layer) Image() *image.RGBA {
	min, max := l.stl.Min, l.stl.Max
	bounds := image.Rect(round(min.X*drawfactor), round(min.Y*drawfactor), round(max.X*drawfactor)+1, round(max.Y*drawfactor)+1)
	img := image.NewRGBA(bounds)
	for _, s := range l.perimeters {
		drawLine(img, perimeterColor, s)
	}
	for _, s := range l.infill {
		drawLine(img, infillColor, s)
	}
	return img
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
