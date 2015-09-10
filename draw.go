package slice

import (
	"image"
	"image/draw"
)

const drawfactor = 20

func (l *Layer) Draw() *image.RGBA {
	min, max := l.stl.Min, l.stl.Max
	bounds := image.Rect(round(min.X*drawfactor), round(min.Y*drawfactor), round(max.X*drawfactor)+1, round(max.Y*drawfactor)+1)
	img := image.NewRGBA(bounds)
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
	for _, s := range l.perimeters {
		drawLine(img, s)
	}
	for _, s := range l.infill {
		drawLine(img, s)
	}
	return img
}

func drawLine(img *image.RGBA, seg *segment) {
	x0, y0, x1, y1 := seg.end1.X*drawfactor, seg.end1.Y*drawfactor, seg.end2.X*drawfactor, seg.end2.Y*drawfactor
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := y1 - y0
	if abs(dx) < 0.5 {
		drawLineVert(img, seg)
		return
	}
	var err float64
	slope := abs(dy / dx)
	y := round(y0)
	yDir := sign(y1 - y0)
	for x := round(x0); x <= round(x1); x++ {
		img.Set(x, y, image.Black)
		err = err + slope
		for err >= 0.5 && !exceeded(float64(y), y1, yDir) {
			img.Set(x, y, image.Black)
			y += yDir
			err -= 1.0
		}
	}
}

func drawLineVert(img *image.RGBA, seg *segment) {
	x := round(seg.end1.X * drawfactor)
	y0, y1 := round(seg.end1.Y*drawfactor), round(seg.end2.Y*drawfactor)
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	for y := y0; y < y1; y++ {
		img.Set(x, y, image.Black)
	}
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

func abs(v float64) float64 {
	if v < 0 {
		return -1 * v
	}
	return v
}

func sign(v float64) int {
	if v < 0 {
		return -1
	}
	return 1
}

func exceeded(from, to float64, dir int) bool {
	if dir < 0 {
		return to > from
	}
	return from > to
}
