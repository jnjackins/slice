package slice

import "fmt"

type Layer struct {
	facets   []*Facet
	segments []segment
}

//TODO: linked list?
//TODO: less brute force
//TODO: sort by lowZ and stop when lowZ > z
func (s *STL) mkLayer(z float64) *Layer {
	facets := make([]*Facet, 0)
	for _, f := range s.facets {
		if f.lowZ <= z && f.highZ >= z {
			facets = append(facets, f)
		}
	}

	segments := make([]segment, 0, len(facets))
	zs := segment{}
	for _, f := range facets {
		s := sliceFacet(f, z)
		if s != zs {
			segments = append(segments, s)
		}
	}
	dprintf("layer z=%0.3f: %d facets / %d segments", z, len(facets), len(segments))

	l := &Layer{
		facets:   facets,
		segments: segments,
	}
	return l
}

type segment struct {
	end1, end2 Vertex
}

func (s segment) String() string {
	return fmt.Sprintf("%v-%v", s.end1, s.end2)
}
