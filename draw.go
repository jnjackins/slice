package slice

import (
	"image"
	"image/draw"
)

func (l *Layer) Draw() *image.RGBA {
	min, max := l.stl.min, l.stl.max
	bounds := image.Rect(round(min.x), round(min.y), round(max.x)+1, round(max.y)+1)
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
	x0, y0, x1, y1 := seg.end1.x, seg.end1.y, seg.end2.x, seg.end2.y
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	deltaX := x1 - x0
	deltaY := y1 - y0
	if deltaX == 0 {
		drawLineVert(img, seg)
		return
	}
	var err float64
	deltaErr := abs(deltaY / deltaX)
	y := round(y0)
	yDir := sign(y1 - y0)
	for x := round(x0); x <= round(x1); x++ {
		img.Set(x, y, image.Black)
		err = err + deltaErr
		for err >= 0.5 {
			img.Set(x, y, image.Black)
			y += yDir
			err -= 1.0
		}
	}
}

func drawLineVert(img *image.RGBA, seg *segment) {
	x := round(seg.end1.x)
	y0, y1 := round(seg.end1.y), round(seg.end2.y)
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	for y := y0; y <= y1; y++ {
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
