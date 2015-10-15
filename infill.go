package slice

// Infill algorithm (for one layer):
//
// - for each segment s, set s.first and s.second to some permutation of s.from and s.to in relation to a line sortLine with angle InfillAngle
// - sort all segments into a list l1, and all segments into a list l2 by distance of s.first (l1) and s.second (l2) from sortLine
// - use the same process for two lists l1Perp and l2Perp, sorted using InfillAngle rotated by 90° (used for "turns" in the infill)
// - all segments are initially white
// loop1: until all segments are marked grey:
//   - set dot to vector v, defined by a point (segment end closest to (0,0) whose segment is not grey) and an angle InfillAngle
//   - shift dot and v down-right so that dot is aligned with the next good InfillSpacing value (TODO - how to calculate this?)
//     -- if such a value does not exist on the segment, mark the segment grey and go back to loop1
//   loop2: until there are no intesections with v:
//     - lookup (binary search) segments in l1 with dist(first, v) < 0 (matches1)
//     - lookup segments in l2 with with dist(second, v) > 0 (matches2)
//     - find all segments in both matches1 and matches2 - these are segments that intersect with v
//     - calculate intersection points with those segments, find closest
//     - advance dot to the intersection, move v to dot and rotate v by 90° (or -90°)
//     - mark the segment grey
//     - repeat the turn-advance process until we are back near the starting point and facing the same direction as InfillAngle again
//
// Memory requirements: (16 bytes per Vertex2) * (2 Vertex2s per segment) * (4 lists) = 128 bytes per segment

import (
	"math"
	"sort"
)

const (
	markWhite = iota
	markGrey
)

func (l *Layer) genInfill(cfg Config) {
	dprintf("generating infill for layer %d...", l.n)
	l.infill = make([]*segment, 0)

	infillAngle := cfg.InfillAngle * math.Pi / 180.0
	if l.n%2 == 1 {
		infillAngle += math.Pi / 2.0
	}

	l1, l2 := l.sortSegments(infillAngle)
	_ = l2
	dot := l1[0].first
	dprintf("dot: %v", dot)
	cast := newLine(dot, infillAngle)
	dprintf("cast: %v", cast)
	x1 := l.stl.Min.X
	y1 := cast.m*x1 + cast.b
	x2 := l.stl.Max.X
	y2 := cast.m*x2 + cast.b
	l.infill = append(l.infill, &segment{from: Vertex2{x1, y1}, to: Vertex2{x2, y2}})
}

type byDist struct {
	data []*segment
	end  int
}

func (a byDist) Len() int      { return len(a.data) }
func (a byDist) Swap(i, j int) { a.data[i], a.data[j] = a.data[j], a.data[i] }
func (a byDist) Less(i, j int) bool {
	if a.end == 1 {
		return a.data[i].dfirst < a.data[j].dfirst
	} else if a.end == 2 {
		return a.data[i].dsecond < a.data[j].dsecond
	} else {
		panic("byDist.Less: invalid end")
	}
}

func (l *Layer) sortSegments(angle float64) ([]*segment, []*segment) {
	// ensure that all segments are on one side of the sorting line by placing it
	// on the top-left corner if the slope is positive, and the top-right corner if
	// the slope is negative.
	var origin Vertex2
	if math.Tan(angle) < 0 {
		origin = Vertex2{l.stl.Min.X, l.stl.Min.Y}
	} else {
		origin = Vertex2{l.stl.Max.X, l.stl.Min.Y}
	}
	sortLine := newLine(origin, angle)
	for _, s := range l.perimeters {
		d1 := sortLine.dist(s.from)
		d2 := sortLine.dist(s.to)
		if d2 > d1 {
			s.first = s.from
			s.dfirst = d1
			s.second = s.to
			s.dsecond = d2
		} else {
			s.first = s.to
			s.dfirst = d2
			s.second = s.from
			s.dsecond = d1
		}
	}
	l1 := make([]*segment, len(l.perimeters))
	l2 := make([]*segment, len(l.perimeters))
	copy(l1, l.perimeters)
	copy(l2, l.perimeters)
	sort.Sort(byDist{l1, 1})
	sort.Sort(byDist{l2, 2})
	return l1, l2
}
