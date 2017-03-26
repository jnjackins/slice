package vector

import (
	"fmt"
	"math"
)

type V3 struct {
	X, Y, Z float64
}

func (v1 V3) Sub(v2 V3) V3 {
	return V3{X: v1.X - v2.X, Y: v1.Y - v2.Y, Z: v1.Z - v2.Z}
}

func (v1 V3) Add(v2 V3) V3 {
	return V3{X: v1.X + v2.X, Y: v1.Y + v2.Y, Z: v1.Z + v2.Z}
}

func (v V3) Mul(d float64) V3 {
	return V3{X: v.X * d, Y: v.Y * d, Z: v.Z * d}
}

func (v V3) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z))
}
func (v V3) Norm() V3 {
	length := v.Length()
	return V3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

func (v V3) String() string {
	return fmt.Sprintf("(%.1f, %.1f, %.1f)", v.X, v.Y, v.Z)
}
