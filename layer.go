package slice

import "fmt"

type Layer struct {
	stl        *STL
	facets     []*facet
	perimeters []*segment
	infill     []*segment
}

type segment struct {
	end1, end2 vertex
}

func (s *segment) String() string {
	return fmt.Sprintf("%v-%v", s.end1, s.end2)
}

func (s *STL) sliceLayer(z float64) *Layer {
	// find the facets which interect this layer
	facets := make([]*facet, 0)
	for _, f := range s.facets {
		if f.lowZ <= z && f.highZ >= z {
			facets = append(facets, f)
		}
	}

	// slice each facet to find the perimeters
	segments := make([]*segment, 0, len(facets))
	zs := segment{}
	for _, f := range facets {
		s := sliceFacet(f, z)
		if s != zs {
			segments = append(segments, &s)
		}
	}
	dprintf("layer z=%0.3f: %d facets / %d perimeter segments", z, len(facets), len(segments))

	l := &Layer{
		stl:        s,
		facets:     facets,
		perimeters: segments,
	}

	l.genInfill()

	return l
}
