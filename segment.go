package slice

import (
	"fmt"
	"math"

	"sigint.ca/slice/internal/vector"
)

var errNoIntersections = fmt.Errorf("no intersections")
var errOutOfBounds = fmt.Errorf("line out of bounds")

// segments can represent both perimeter and infill lines
type segment struct {
	from, to Vertex2   // ordered so that gcode movements are from "from" to "to"
	normal   vector.V2 // points to the inside of the solid
	line     *line     // hold on to the slope and y-intercept of the line once calculated
	visited  bool      // used to keep track of regions that still need to be infilled
}

func (s *segment) String() string {
	return fmt.Sprintf("%v-%v", s.from, s.to)
}

func (s *segment) getLine() *line {
	if s.line != nil {
		dprintf("returning cached line for %v", s)
		return s.line
	}
	div := s.to.X - s.from.X
	var slope float64
	if div == 0 {
		slope = math.Inf(+1)
		s.line = &line{m: slope, b: s.to.X}
	} else {
		slope = (s.to.Y - s.from.Y) / div
		s.line = &line{m: slope, b: s.from.Y - slope*s.from.X}
	}
	return s.line
}

func (s *segment) shiftBy(v vector.V2) *segment {
	dprintf("shifting by %v", v)
	s.from = Vertex2(vector.V2(s.from).Add(v))
	s.to = Vertex2(vector.V2(s.to).Add(v))
	s.line = nil
	return s
}

// getIntersections returns a list of segments in target that intersect with ray, as well
// as a list of the corresponding intersection points.
func (ray *segment) getIntersections(target []*segment) ([]*segment, []Vertex2) {
	intersecting := make([]*segment, 0)
	points := make([]Vertex2, 0)
	for _, s := range target {
		// eliminate cases where the segments do not have overlapping X coordinates
		if math.Max(ray.from.X, ray.to.X) < math.Min(s.from.X, s.to.X) {
			continue
		}
		if math.Min(ray.from.X, ray.to.X) > math.Max(s.from.X, s.to.X) {
			continue
		}

		l1, l2 := ray.getLine(), s.getLine()
		if approxEquals(l1.m, l2.m, 0.000001) {
			// l1 and l2 are parallel
			continue
		}

		// calculate point of intersection
		var v Vertex2
		if math.IsInf(l1.m, 0) {
			// ray is vertical
			v.X = ray.from.X
			v.Y = l2.m*v.X + l2.b
			if !inRange(v.Y, s.from.Y, s.to.Y) {
				continue
			}
		} else if math.IsInf(l2.m, 0) {
			// s is vertical
			v.X = s.from.X
			v.Y = l1.m*v.X + l1.b
			if !inRange(v.Y, s.from.Y, s.to.Y) {
				continue
			}
		} else {
			v.X = (l2.b - l1.b) / (l1.m - l2.m)
			v.Y = l1.m*v.X + l1.b // doesn't matter which line we use in this case
			if !inRange(v.X, s.from.X, s.to.X) {
				continue
			}
		}

		intersecting = append(intersecting, s)
		points = append(points, v)
	}
	return intersecting, points
}

// checkDomain returns true if v is within the domain x1..x2, or false otherwise
func inRange(test, v1, v2 float64) bool {
	return test >= math.Min(v1, v2)-0.000001 && test <= math.Max(v1, v2)+0.000001
}

type line struct {
	m float64 // slope of line (calculated)
	b float64 // y intercept (calculated) (or x value if m is infinite)
}

func (l *line) String() string {
	if math.IsInf(l.m, 0) {
		return fmt.Sprintf("vertical line at x=%v", l.b)
	}
	return fmt.Sprintf("(y=%.1fx+%.1f", l.m, l.b)
}

func lineFromAngle(origin Vertex2, angle float64) *line {
	slope := math.Tan(angle)
	return &line{m: slope, b: origin.Y - slope*origin.X}
}

func (l1 *line) intersect(l2 *line) (Vertex2, error) {
	var x, y float64
	if math.IsInf(l1.m, 0) && math.IsInf(l2.m, 0) {
		return Vertex2{}, errNoIntersections
	}
	if math.IsInf(l1.m, 0) {
		x = l1.b
		y = l2.m*x + l2.b // l2 is not vertical
	} else if math.IsInf(l2.m, 0) {
		x = l2.b
		y = l1.m*x + l1.b // l1 is not vertical
	} else {
		div := l1.m - l2.m
		if div == 0 {
			return Vertex2{}, errNoIntersections
		}
		x = (l2.b - l1.b) / div
		y = l1.m*x + l1.b // doesn't matter which line is used here
	}

	return Vertex2{X: x, Y: y}, nil
}

// bound returns a segment representing the line bounded by the rectangle (min-max)
func (l *line) bound(min, max Vertex2) (*segment, error) {
	left := &segment{from: min, to: Vertex2{X: min.X, Y: max.Y}}
	right := &segment{from: Vertex2{X: max.X, Y: min.Y}, to: max}
	top := &segment{from: min, to: Vertex2{X: max.X, Y: min.Y}}
	bottom := &segment{from: Vertex2{X: min.X, Y: max.Y}, to: max}

	ends := make([]Vertex2, 0, 2)

	v1, err := left.getLine().intersect(l)
	if err == nil && inRange(v1.Y, left.from.Y, left.to.Y) {
		ends = append(ends, v1)
	}
	v2, err := right.getLine().intersect(l)
	if err == nil && inRange(v2.Y, right.from.Y, right.to.Y) {
		ends = append(ends, v2)
	}
	v3, err := top.getLine().intersect(l)
	if err == nil && inRange(v3.X, top.from.X, top.to.X) && !v3.touches(v1) && !v3.touches(v2) {
		ends = append(ends, v3)
	}
	v4, err := bottom.getLine().intersect(l)
	if err == nil && inRange(v4.X, bottom.from.X, bottom.to.X) && !v4.touches(v1) && !v4.touches(v2) {
		ends = append(ends, v4)
	}

	if len(ends) != 2 {
		return nil, errOutOfBounds
	}

	return &segment{
		from: ends[0],
		to:   ends[1],
		line: l,
	}, nil
}
