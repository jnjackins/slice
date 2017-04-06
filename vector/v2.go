package vector

import (
	"fmt"
	"math"
)

type V2 struct {
	X, Y float64
}

func (v1 V2) Sub(v2 V2) V2 {
	return V2{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 V2) Add(v2 V2) V2 {
	return V2{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v V2) Mul(d float64) V2 {
	return V2{X: v.X * d, Y: v.Y * d}
}

func (v V2) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}
func (v V2) Normalize() V2 {
	length := v.Length()
	return V2{X: v.X / length, Y: v.Y / length}
}

func (v V2) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", v.X, v.Y)
}
