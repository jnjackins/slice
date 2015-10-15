package slice

import "math"

type line struct {
	origin Vertex2 // a point on the line
	angle  float64 // angle of the line, relative to the line y = origin.Y
	m      float64 // slope of line (calculated)
	b      float64 // y intercept (calculated)
}

func newLine(origin Vertex2, angle float64) line {
	slope := math.Tan(angle)
	return line{origin: origin, angle: angle, m: slope, b: origin.Y - slope*origin.X}
}

func (l line) dist(v Vertex2) float64 {
	m := -1 / l.m    // slope of line perpendicular to m (l2)
	b := v.Y - m*v.X // y-intercept of l2
	x := (b - l.b) / (l.m - m)
	y := m*x + b
	d := math.Sqrt(math.Pow(y-v.Y, 2) + math.Pow(x-v.X, 2))
	return d
}
