package slice

import (
	"fmt"

	"sigint.ca/slice/vector"
)

type Vertex2 struct {
	X, Y float64
}

func (v1 Vertex2) touches(v2 Vertex2) bool {
	return approxEquals(v1.X, v2.X, 0.005) && approxEquals(v1.Y, v2.Y, 0.005)
}

func (v1 Vertex2) distFrom(v2 Vertex2) float64 {
	v := vector.V2(v2).Sub(vector.V2(v1))
	return v.Length()
}

func (v Vertex2) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f)", v.X, v.Y)
}
