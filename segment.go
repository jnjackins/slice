package slice

import (
	"fmt"
	"math"

	"sigint.ca/slice/internal/vector"
)

const (
	markWhite = iota
	markGrey
)

// segments can represent both perimeter and infill lines
type segment struct {
	from, to Vertex2 // ordered so that gcode movements are from "from" to "to"

	line *line // hold on to the slope and y-intercept of the line once calculated
	mark int   // used to keep track of regions that still need to be infilled
}

func (s *segment) String() string {
	return fmt.Sprintf("%v-%v", s.from, s.to)
}

func (s *segment) getLine() *line {
	if s.line != nil {
		return s.line
	}
	div := s.to.X - s.from.X
	if div == 0 {
		wprintf("lineFromSegment: division by 0")
	}
	slope := (s.to.Y - s.from.Y) / div
	s.line = &line{m: slope, b: s.from.Y - slope*s.from.X}
	return s.line
}

func (s *segment) shiftBy(v vector.V2) {
	dprintf("shifting by %v", v)
	s.from = Vertex2(vector.V2(s.from).Add(v))
	s.to = Vertex2(vector.V2(s.to).Add(v))
}

// getIntersections returns a list of segments in target that intersect with ray, as well
// as a list of the corresponding intersection points.
func (ray *segment) getIntersections(target []*segment) ([]*segment, []Vertex2) {
	intersecting := make([]*segment, 0)
	points := make([]Vertex2, 0)
	for _, s := range target {
		// eliminate cases where the segments do not have overlapping X coordinates
		if math.Min(ray.from.X, ray.to.X) > math.Max(s.from.X, s.to.X) {
			continue
		}
		if math.Max(ray.from.X, ray.to.X) < math.Min(s.from.X, s.to.X) {
			continue
		}

		// eliminate cases where the segments are parallel
		l1, l2 := ray.getLine(), s.getLine()
		if approxEquals(l1.m, l2.m) {
			continue
		}

		// calculate point of intersection
		x := (l2.b - l1.b) / (l1.m - l2.m) // non-zero divisor (verified slopes are not equal above)
		y := l1.m*x + l1.b

		// test that the intersection is within the domain of the line segments
		if (x < math.Max(math.Min(ray.from.X, ray.to.X), math.Min(s.from.X, s.to.X))) ||
			(x > math.Min(math.Max(ray.from.X, ray.to.X), math.Max(s.from.X, s.to.X))) {
			continue
		}
		intersecting = append(intersecting, s)
		points = append(points, Vertex2{X: x, Y: y})
	}
	return intersecting, points
}

type line struct {
	m float64 // slope of line (calculated)
	b float64 // y intercept (calculated)
}

func (l line) String() string {
	return fmt.Sprintf("(y=%.1fx+%.1f", l.m, l.b)
}

func lineFromAngle(origin Vertex2, angle float64) line {
	slope := math.Tan(angle)
	return line{m: slope, b: origin.Y - slope*origin.X}
}

func (l1 line) intersectionPoint(l2 line) Vertex2 {
	div := l1.m - l2.m
	if div == 0 {
		wprintf("intersectionPoint: division by 0")
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
