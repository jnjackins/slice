package slice

import (
	"fmt"
	"math"
)

type vector Vertex2

func (v1 vector) sub(v2 vector) vector {
	return vector{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 vector) add(v2 vector) vector {
	return vector{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v vector) mul(d float64) vector {
	return vector{X: v.X * d, Y: v.Y * d}
}

func (v vector) length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}
func (v vector) norm() vector {
	length := v.length()
	return vector{X: v.X / length, Y: v.Y / length}
}

func (v vector) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", v.X, v.Y)
}
