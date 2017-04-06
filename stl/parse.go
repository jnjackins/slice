package stl

import (
	"bufio"
	"fmt"
	"io"

	"sigint.ca/slice/vector"
)

type facetReader interface {
	readFacets() ([]Facet, error)
}

type Facet struct {
	Vertices [3]vector.V3
	Normal   vector.V3
}

func (f Facet) String() string {
	return fmt.Sprintf("n=%v v=%v\n", f.Normal, f.Vertices)
}

// Parse parses a new STL from an io.Reader.
func Parse(r io.Reader) (*Solid, error) {
	bufr := bufio.NewReader(r)
	top, err := bufr.Peek(6)
	if err != nil {
		return nil, err
	}
	var fr facetReader
	if string(top) == "solid " {
		fr = &asciiReader{bufr}
	} else {
		fr = &binaryReader{bufr}
	}
	facets, err := fr.readFacets()
	if err != nil {
		return nil, err
	}

	s := &Solid{Facets: facets}
	s.updateBounds()

	return s, nil
}

func (f *Facet) calculateNormal() {
	v := vector.V3(f.Vertices[1]).Sub(vector.V3(f.Vertices[0]))
	w := vector.V3(f.Vertices[2]).Sub(vector.V3(f.Vertices[0]))
	n := vector.V3{
		X: (v.Y * w.Z) - (v.Z * w.Y),
		Y: (v.Z * w.X) - (v.X * w.Z),
		Z: (v.X * w.Y) - (v.Y * w.X),
	}
	f.Normal = n.Normalize()
}
