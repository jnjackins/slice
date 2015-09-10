package slice

import (
	"image"

	"sigint.ca/graphics/primitive"
)

const drawfactor = 20

func (l *Layer) Image() *image.RGBA {
	min, max := l.stl.Min, l.stl.Max
	bounds := image.Rect(round(min.X*drawfactor), round(min.Y*drawfactor), round(max.X*drawfactor)+1, round(max.Y*drawfactor)+1)
	img := image.NewRGBA(bounds)
	for _, s := range l.perimeters {
		drawLine(img, s)
	}
	for _, s := range l.infill {
		drawLine(img, s)
	}
	return img
}

func drawLine(img *image.RGBA, seg *segment) {
	p1 := image.Pt(round(seg.end1.X*drawfactor), round(seg.end1.Y*drawfactor))
	p2 := image.Pt(round(seg.end2.X*drawfactor), round(seg.end2.Y*drawfactor))
	primitive.Line(img, p1, p2)
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
