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
	perimeterColor = color.Black
	infillColor    = color.RGBA{R: 0xFF, A: 0xFF}
	normalColor    = color.RGBA{G: 0xFF, A: 0xFF}
)

func (l *Layer) Draw(dst draw.Image) {
	// scale to window size
	min, max := l.stl.Bounds()
	min2 := Vertex2{X: min.X, Y: min.Y}
	max2 := Vertex2{X: max.X, Y: max.Y}
	srcDx := max.X - min.X
	srcDy := max.Y - min.Y
	dr := dst.Bounds()
	dstDx := float64(dr.Dx())
	dstDy := float64(dr.Dy())
	var scaleFactor float64
	if dstDx/srcDx < dstDy/srcDy {
		scaleFactor = dstDx / srcDx
	} else {
		scaleFactor = dstDy / srcDy
	}

	// show stl bounds
	r := image.Rect(0, 0, int((max2.X-min2.X)*scaleFactor), int((max2.Y-min2.Y)*scaleFactor))
	draw.Draw(dst, r, image.White, image.ZP, draw.Src)

	// draw the layer
	for _, region := range l.Regions() {
		for i, s := range region.Exterior {
			drawNumber(dst, s.From, min2, fmt.Sprintf("%d", i), scaleFactor)
			drawSegment(dst, s, min2, perimeterColor, scaleFactor)
		}
		for _, p := range region.Interiors {
			for _, s := range p {
				drawSegment(dst, s, min2, perimeterColor, scaleFactor)
			}
		}
		for _, s := range region.Infill {
			drawSegment(dst, s, min2, infillColor, scaleFactor)
		}
	}
}

func drawSegment(dst draw.Image, s *Segment, min Vertex2, c color.Color, scaleFactor float64) {
	px1 := v2pixel(s.From, scaleFactor).Sub(v2pixel(min, scaleFactor))
	px2 := v2pixel(s.To, scaleFactor).Sub(v2pixel(min, scaleFactor))
	primitive.Circle(dst, c, px1, 2)
	primitive.Circle(dst, c, px1, 2)
	primitive.Line(dst, c, px1, px2)

	// draw normal
	n1 := image.Pt((px1.X+px2.X)/2, (px1.Y+px2.Y)/2)
	normal := image.Pt(int(s.Normal.X*scaleFactor*0.1), int(s.Normal.Y*scaleFactor*0.1))
	n2 := n1.Add(normal)
	primitive.Line(dst, normalColor, n1, n2)
}

func drawNumber(dst draw.Image, pt, min Vertex2, number string, scaleFactor float64) {
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
