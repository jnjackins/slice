package slice

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"sigint.ca/graphics/primitive"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	PerimeterColor = color.Black
	InfillColor    = color.RGBA{R: 0xFF, A: 0xFF}
)

func (l *Layer) Draw(dst draw.Image) {
	var scaleFactor float64
	min, max := l.stl.Bounds()
	min2 := Vertex2{X: min.X, Y: min.Y}
	srcDx := max.X - min.X
	srcDy := max.Y - min.Y
	r := dst.Bounds()
	dstDx := float64(r.Dx())
	dstDy := float64(r.Dy())
	if dstDx-srcDx < dstDy-srcDy {
		scaleFactor = dstDx / srcDx
	} else {
		scaleFactor = dstDy / srcDy
	}

	for i, region := range l.Regions() {
		drawPerimeterNumber(dst, region.Exterior[0].From, min2, fmt.Sprintf("%d", i), scaleFactor)
		for _, s := range region.Exterior {
			drawLine(dst, s.From, s.To, min2, PerimeterColor, scaleFactor)
		}
		for j, p := range region.Interiors {
			drawPerimeterNumber(dst, p[0].From, min2, fmt.Sprintf("%d-%d", i, j), scaleFactor)
			for _, s := range p {
				drawLine(dst, s.From, s.To, min2, PerimeterColor, scaleFactor)
			}
		}
		for _, s := range region.Infill {
			drawLine(dst, s.From, s.To, min2, InfillColor, scaleFactor)
		}
	}
}

func drawLine(dst draw.Image, p1, p2, min Vertex2, c color.Color, scaleFactor float64) {
	px1 := v2pixel(p1, scaleFactor).Sub(v2pixel(min, scaleFactor))
	px2 := v2pixel(p2, scaleFactor).Sub(v2pixel(min, scaleFactor))
	primitive.Line(dst, c, px1, px2)
}

func drawPerimeterNumber(dst draw.Image, pt, min Vertex2, number string, scaleFactor float64) {
	px := v2pixel(pt, scaleFactor).Sub(v2pixel(min, scaleFactor))
	dot := fixed.P(px.X, px.Y)
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA{B: 0xFF, A: 0xFF}),
		Face: basicfont.Face7x13,
		Dot:  dot,
	}
	d.DrawString(number)
}

func v2pixel(v Vertex2, scaleFactor float64) image.Point {
	return image.Pt(round(v.X*scaleFactor), round(v.Y*scaleFactor))
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
