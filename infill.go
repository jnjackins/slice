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
	l.debug = make([]*segment, 0)

	infillAngle := cfg.InfillAngle * math.Pi / 180.0
	if l.n%2 == 1 {
		infillAngle += math.Pi / 2.0
	}

	var origin Vertex2 // infill proceeds away from origin
	var infillDir int
	slope := math.Tan(infillAngle)
	if slope < 0 {
		origin = Vertex2{l.stl.Min.X, l.stl.Min.Y}
		infillDir = -1
	} else if slope > 0 {
		origin = Vertex2{l.stl.Max.X, l.stl.Min.Y}
		infillDir = 1
	} else {
		infillDir = 0
	}

	// sort segments into two lists, by distance of their endpoints from a
	// sorting line with angle infillAngle. One list for each segment endpoint.
	l1, l2 := l.sortSegments(infillAngle, origin)

	// Find a starting point. Start from the first end of the least distant
	// segment (as determined by the sorting), shifted by cfg.LineWidth.
	// Create a line at that point, with angle infillAngle.
	// If there are at least 2 segments that intersect with that line,
	// set dot the intersection point on the line that is closest to one
	// end of the line (the "top" end). Otherwise, TODO
	dprintf("looking for starting point")
	castLine := lineFromAngle(l1[0].first, infillAngle)
	from := Vertex2{l.stl.Min.X, castLine.m*l.stl.Min.X + castLine.b}
	to := Vertex2{l.stl.Max.X, castLine.m*l.stl.Max.X + castLine.b}
	// cast from the top towards the bottom
	// TODO: don't assume that cast isn't horizontal
	if to.Y < from.Y {
		to, from = from, to
	}
	cast := &segment{from: from, to: to}
	shiftAngle := infillAngle + math.Pi/2.0
	shiftLine := lineFromAngle(origin, shiftAngle)

	// shift is vector representing the direction that we need to shift the
	// cast line
	v := vector(shiftLine.intersectionPoint(castLine)).sub(vector(origin))
	shift := v.norm().mul(cfg.LineWidth)

	// shift the cast line inwards by cfg.LineWidth
	dprintf("shifting by %v (%.1f°)", shift, shiftAngle*180/math.Pi)
	cast.shiftBy(shift)
	castLine = lineFromSegment(cast)

	dprintf("trying cast=%v (castLine=%v)", cast, castLine)

	intersections := l.getIntersections(cast, infillDir, l1, l2)

	// TODO: what do we do if there aren't 2 intersections?
	n := len(intersections)
	dprintf("%d intersections", n)
	if n >= 2 {
		// get exact intersection points
		points := make([]Vertex2, len(intersections))
		for i, s := range intersections {
			points[i] = lineFromSegment(s).intersectionPoint(castLine)
		}

		// use the first two points
		sort.Sort(verticesByDist{points, cast.from})
		s := &segment{from: points[0], to: points[1]}

		dprintf("adding infill segment: %v", s)
		l.infill = append(l.infill, s)
	}
	//l.debug = append(l.debug, cast)
}

type verticesByDist struct {
	points []Vertex2
	from   Vertex2
}

func (a verticesByDist) Len() int      { return len(a.points) }
func (a verticesByDist) Swap(i, j int) { a.points[i], a.points[j] = a.points[j], a.points[i] }
func (a verticesByDist) Less(i, j int) bool {
	return a.points[i].distFrom(a.from) < a.points[j].distFrom(a.from)
}

type segmentsByDist struct {
	data []*segment
	end  int
}

func (a segmentsByDist) Len() int      { return len(a.data) }
func (a segmentsByDist) Swap(i, j int) { a.data[i], a.data[j] = a.data[j], a.data[i] }
func (a segmentsByDist) Less(i, j int) bool {
	if a.end == 1 {
		return a.data[i].dfirst < a.data[j].dfirst
	} else if a.end == 2 {
		return a.data[i].dsecond < a.data[j].dsecond
	} else {
		panic("segmentsByDist.Less: invalid end")
	}
}

// sortSegments returns two lists of sorted segments. The lists are sorted
// by the segment's distance from the sorting line, one list for each segment end.
func (l *Layer) sortSegments(angle float64, origin Vertex2) (l1, l2 []*segment) {
	sortLine := lineFromAngle(origin, angle)
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
	l1 = make([]*segment, len(l.perimeters))
	l2 = make([]*segment, len(l.perimeters))
	copy(l1, l.perimeters)
	copy(l2, l.perimeters)
	sort.Sort(segmentsByDist{l1, 1})
	sort.Sort(segmentsByDist{l2, 2})
	return l1, l2
}

func (l *Layer) getIntersections(cast *segment, infillDir int, l1, l2 []*segment) []*segment {
	i := sort.Search(len(l1), func(i int) bool {
		return infillDir*checkSide(cast, l1[i].first) >= 0
	})
	matches1 := l1[:i]
	j := sort.Search(len(l2), func(i int) bool {
		return infillDir*checkSide(cast, l2[i].second) >= 0
	})
	matches2 := l2[j:]

	matchMap := make(map[*segment]int)
	intersections := make([]*segment, 0)
	for _, s := range matches1 {
		matchMap[s]++
	}

	for _, s := range matches2 {
		if _, ok := matchMap[s]; ok {
			dprintf("%v intersects with %v", s, cast)
			intersections = append(intersections, s)
		}
	}

	return intersections
}

func traverse(dot Vertex2, s *segment, d float64) Vertex2 {
	n := vector(s.second).sub(vector(s.first)).norm()
	dprintf("normal of %v: %v", s, n)
	v := vector(dot).add(n.mul(d))
	return Vertex2(v)
}

// checkSide returns -1, +1, or 0 if p is on one side of s, the other, or directly on s.
func checkSide(s *segment, p Vertex2) int {
	position := sign((s.to.X-s.from.X)*(p.Y-s.from.Y) - (s.to.Y-s.from.Y)*(p.X-s.from.X))
	return position
}

func sign(v float64) int {
	if v < 0.0 {
		return -1
	}
	if v > 0.0 {
		return 1
	}
	return 0
}
