package slice

import (
	"fmt"
	"image"

	"sigint.ca/slice/internal/vector"
)

type Vertex2 struct {
	X, Y float64
}

func (v1 Vertex2) touches(v2 Vertex2) bool {
	return approxEquals(v1.X, v2.X, 0.000001) && approxEquals(v1.Y, v2.Y, 0.000001)
}

func (v1 Vertex2) near(v2 Vertex2) bool {
	return approxEquals(v1.X, v2.X, 0.001) && approxEquals(v1.Y, v2.Y, 0.001)
}

func (v1 Vertex2) distFrom(v2 Vertex2) float64 {
	v := vector.V2(v2).Sub(vector.V2(v1))
	return v.Length()
}

func (v Vertex2) pt() image.Point {
	return image.Pt(round(v.X*drawfactor), round(v.Y*drawfactor))
}

func (v Vertex2) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f)", v.X, v.Y)
}
