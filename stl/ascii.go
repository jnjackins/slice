package stl

import (
	"bufio"
	"bytes"
	"fmt"

	"sigint.ca/slice/vector"
)

type asciiReader struct {
	r *bufio.Reader
}

func (p asciiReader) readFacets() ([]Facet, error) {
	// discard header
	p.r.ReadLine()

	facets := make([]Facet, 0)
	for {
		p.r.ReadLine() // discard normal
		p.r.ReadLine() // discard "outer loop"

		f, err := p.readFacet()
		if err != nil {
			return nil, fmt.Errorf("error decoding ascii STL: %v", err)
		}
		facets = append(facets, f)

		p.r.ReadLine() // discard "endloop"
		p.r.ReadLine() // discard "endfacet"

		nextWord, err := p.r.Peek(8)
		if err != nil {
			return nil, fmt.Errorf("error decoding ascii STL: %v", err)
		}
		if string(nextWord) == "endsolid" {
			break
		}
	}

	return facets, nil
}

func (p asciiReader) readFacet() (Facet, error) {
	var vertices [3]vector.V3
	for vi := range vertices {
		v, err := p.readVertex()
		if err != nil {
			return Facet{}, fmt.Errorf("read facet: %v", err)
		}
		vertices[vi] = v
	}
	f := Facet{Vertices: vertices}
	f.calculateNormal()
	return f, nil
}

func (p asciiReader) readVertex() (vector.V3, error) {
	var x, z, y float32

	s, _, err := p.r.ReadLine()
	if err != nil {
		return vector.V3{}, fmt.Errorf("read vertex: %v", err)
	}
	s = bytes.TrimSpace(s)

	if _, err := fmt.Sscanf(string(s), "vertex %f %f %f\n", &x, &y, &z); err != nil {
		return vector.V3{}, fmt.Errorf("read vertex: %q: %v", s, err)
	}
	v := vector.V3{X: float64(x), Y: float64(y), Z: float64(z)}
	return v, nil
}
