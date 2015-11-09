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

const drawfactor = 10

var (
	perimeterColor = color.Black
	infillColor    = color.RGBA{R: 0xFF, A: 0xFF}
	debugColor     = color.RGBA{G: 0xFF, A: 0xFF}
)

// Bounds returns an image.Rectangle for use with Draw.
func (l *Layer) Bounds() image.Rectangle {
	x1, y1 := int(l.stl.Min.X*drawfactor-20), int(l.stl.Min.Y*drawfactor-20)
	x2, y2 := int(l.stl.Max.X*drawfactor+20.5), int(l.stl.Max.Y*drawfactor+20.5)
	return image.Rect(x1, y1, x2, y2)
}

// Draw draws an image representation of the layer onto dst. Use Bounds to find
// the minimum size that dst should be to contain the entire image.
func (l *Layer) Draw(dst draw.Image) {
	for i, solid := range l.solids {
		drawPerimeterNumber(dst, solid.exterior[0].from.pt(), fmt.Sprintf("%d", i))
		for _, s := range solid.exterior {
			drawLine(dst, s.from, s.to, perimeterColor)
		}
		for j, p := range solid.interiors {
			drawPerimeterNumber(dst, p[0].from.pt(), fmt.Sprintf("%d-%d", i, j))
			for _, s := range p {
				drawLine(dst, s.from, s.to, perimeterColor)
			}
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
	primitive.Line(dst, c, p1.pt(), p2.pt())
}

func drawPerimeterNumber(dst draw.Image, pt image.Point, number string) {
	dot := fixed.P(pt.X, pt.Y)
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA{B: 0xFF, A: 0xFF}),
		Face: basicfont.Face7x13,
		Dot:  dot,
	}
	d.DrawString(number)
}
