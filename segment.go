package slice

import (
	"fmt"
	"math"

	"sigint.ca/slice/vector"
)

var errNoIntersections = fmt.Errorf("no intersections")
var errOutOfBounds = fmt.Errorf("line out of bounds")

// Segments can represent both perimeter and infill lines
type Segment struct {
	From, To Vertex2   // ordered so that gcode movements are from "from" to "to"
	Normal   vector.V2 // points to the inside of the solid
	line     *line     // hold on to the slope and y-intercept of the line once calculated
}

func (s *Segment) String() string {
	return fmt.Sprintf("%v-%v", s.From, s.To)
}

func (s *Segment) Length() float64 {
	dx := s.To.X - s.From.Y
	dy := s.To.Y - s.From.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (s *Segment) getLine() *line {
	if s.line != nil {
		return s.line
	}
	div := s.To.X - s.From.X
	var slope float64
	if div == 0 {
		slope = math.Inf(+1)
		s.line = &line{m: slope, b: s.To.X}
	} else {
		slope = (s.To.Y - s.From.Y) / div
		s.line = &line{m: slope, b: s.From.Y - slope*s.From.X}
	}
	return s.line
}

func (s *Segment) ShiftBy(v vector.V2) *Segment {
	s.From = Vertex2(vector.V2(s.From).Add(v))
	s.To = Vertex2(vector.V2(s.To).Add(v))
	s.line = nil
	return s
}

func (s *Segment) intersect(ss *Segment) (Vertex2, bool) {
	if s.From.touches(ss.From) ||
		s.From.touches(ss.To) ||
		s.To.touches(ss.From) ||
		s.To.touches(ss.To) {
		return Vertex2{}, false
	}

	// eliminate cases where the Segments do not have overlapping X coordinates
	if math.Max(s.From.X, s.To.X) < math.Min(ss.From.X, ss.To.X) {
		return Vertex2{}, false
	}
	if math.Min(s.From.X, s.To.X) > math.Max(ss.From.X, ss.To.X) {
		return Vertex2{}, false
	}

	l1, l2 := s.getLine(), ss.getLine()
	if approxEquals(l1.m, l2.m, 0.000001) {
		// l1 and l2 are parallel
		return Vertex2{}, false
	}

	// calculate point of intersection
	// TODO: duplicate of line.intersect?
	var v Vertex2
	if math.IsInf(l1.m, 0) {
		// ray is vertical
		v.X = s.From.X
		v.Y = l2.m*v.X + l2.b
		if !inRange(v.Y, ss.From.Y, ss.To.Y) {
			return Vertex2{}, false
		}
	} else if math.IsInf(l2.m, 0) {
		// s is vertical
		v.X = ss.From.X
		v.Y = l1.m*v.X + l1.b
		if !inRange(v.Y, ss.From.Y, ss.To.Y) {
			return Vertex2{}, false
		}
	} else {
		v.X = (l2.b - l1.b) / (l1.m - l2.m)
		v.Y = l1.m*v.X + l1.b // doesn't matter which line we use in this case
		if !inRange(v.X, ss.From.X, ss.To.X) {
			return Vertex2{}, false
		}
	}
	return v, true
}

// getIntersections returns a list of Segments in target that intersect with ray, as well
// as a list of the corresponding intersection points.
func (ray *Segment) getIntersections(target []*Segment) ([]*Segment, []Vertex2) {
	intersecting := make([]*Segment, 0)
	points := make([]Vertex2, 0)
	for _, s := range target {
		v, ok := ray.intersect(s)
		if !ok {
			continue
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
	if *l1 == *l2 {
		return Vertex2{}, errNoIntersections
	}

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

// bound returns a Segment representing the line bounded by the rectangle (min-max)
func (l *line) bound(min, max Vertex2) (*Segment, error) {
	left := &Segment{From: min, To: Vertex2{X: min.X, Y: max.Y}}
	right := &Segment{From: Vertex2{X: max.X, Y: min.Y}, To: max}
	top := &Segment{From: min, To: Vertex2{X: max.X, Y: min.Y}}
	bottom := &Segment{From: Vertex2{X: min.X, Y: max.Y}, To: max}

	ends := make([]Vertex2, 0, 2)

	v1, err := left.getLine().intersect(l)
	if err == nil && inRange(v1.Y, left.From.Y, left.To.Y) {
		ends = append(ends, v1)
	}
	v2, err := right.getLine().intersect(l)
	if err == nil && inRange(v2.Y, right.From.Y, right.To.Y) {
		ends = append(ends, v2)
	}
	v3, err := top.getLine().intersect(l)
	if err == nil && inRange(v3.X, top.From.X, top.To.X) && !v3.touches(v1) && !v3.touches(v2) {
		ends = append(ends, v3)
	}
	v4, err := bottom.getLine().intersect(l)
	if err == nil && inRange(v4.X, bottom.From.X, bottom.To.X) && !v4.touches(v1) && !v4.touches(v2) {
		ends = append(ends, v4)
	}

	if len(ends) != 2 {
		return nil, errOutOfBounds
	}

	return &Segment{
		From: ends[0],
		To:   ends[1],
		line: l,
	}, nil
}
