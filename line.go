package slice

import (
	"fmt"
	"log"
	"math"
)

type line struct {
	origin Vertex2 // a point on the line
	angle  float64 // angle of the line, relative to the line y = origin.Y
	m      float64 // slope of line (calculated)
	b      float64 // y intercept (calculated)
}

func (l line) String() string {
	return fmt.Sprintf("(y=%.1fx+%.1f", l.m, l.b)
}

func lineFromAngle(origin Vertex2, angle float64) line {
	slope := math.Tan(angle)
	return line{origin: origin, angle: angle, m: slope, b: origin.Y - slope*origin.X}
}

func lineFromSegment(s *segment) line {
	div := s.to.X - s.from.X
	if div == 0 {
		log.Printf("lineFromSegment: warning: division by 0")
	}
	slope := (s.to.Y - s.from.Y) / div
	return line{m: slope, b: s.from.Y - slope*s.from.X}
}

func (l1 line) intersectionPoint(l2 line) Vertex2 {
	div := l1.m - l2.m
	if div == 0 {
		log.Printf("intersectionPoint: warning: division by 0")
	}
	x := (l2.b - l1.b) / div
	y := l2.m*x + l2.b
	return Vertex2{X: x, Y: y}
}

func (l line) dist(v Vertex2) float64 {
	m := -1 / l.m    // slope of line perpendicular to m (l2)
	b := v.Y - m*v.X // y-intercept of l2
	x := (b - l.b) / (l.m - m)
	y := m*x + b
	d := math.Sqrt(math.Pow(y-v.Y, 2) + math.Pow(x-v.X, 2))
	return d
}
