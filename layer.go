package slice

import (
	"fmt"

	"sigint.ca/slice/internal/vector"
)

type Layer struct {
	n          int // layer index
	z          float64
	stl        *STL
	facets     []*facet
	perimeters []*segment
	infill     []*segment
	debug      []*segment // extra lines to draw for debugging purposes
}

type segment struct {
	from, to        Vertex2 // ordered so that gcode movements are from "from" to "to"
	first, second   Vertex2 // ordered relative to some sorting line (only used for sorting and searching)
	dfirst, dsecond float64 // distance of first and second from sorting line
}

func (s *segment) shiftBy(v vector.V2) {
	s.from = Vertex2(vector.V2(s.from).Add(v))
	s.to = Vertex2(vector.V2(s.to).Add(v))
}

func (s *segment) String() string {
	return fmt.Sprintf("%v-%v", s.from, s.to)
}
