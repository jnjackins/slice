package slice

import "fmt"

type Layer struct {
	n          int // layer index
	z          float64
	stl        *STL
	facets     []*facet
	perimeters []*segment
	infill     []*segment
}

type segment struct {
	from, to        Vertex2 // ordered so that gcode movements are from "from" to "to"
	first, second   Vertex2 // ordered relative to some sorting line (changes throughout lifetime)
	dfirst, dsecond float64 // distance of first and second from sorting line
	mark            int     // indicates whether the segment has yet to be visited during infill
}

type Vertex2 struct {
	X, Y float64
}

func (v Vertex2) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f)", v.X, v.Y)
}

func (s *segment) String() string {
	return fmt.Sprintf("%v-%v", s.from, s.to)
}
