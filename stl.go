package slice

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type STL struct {
	Layers   []*Layer
	Min, Max Vertex3

	facets []*facet
}

type facet struct {
	vertices    [3]Vertex3
	lowZ, highZ float64
}

func (f *facet) String() string {
	return fmt.Sprint(f.vertices)
}

type Vertex3 struct {
	X, Y, Z float64
}

func (v Vertex3) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f,%0.3f)", v.X, v.Y, v.Z)
}

// Parse parses a new STL from an io.Reader.
func Parse(r io.Reader) (*STL, error) {
	// test for ascii stl format
	bufr := bufio.NewReader(r)
	top, err := bufr.Peek(6)
	if err != nil {
		return nil, err
	}
	if string(top) == "solid " {
		return nil, fmt.Errorf("ascii format STL not supported")
	}

	// discard text header
	if _, err := bufr.Discard(80); err != nil {
		return nil, fmt.Errorf("error decoding STL: %v", err)
	}

	var nfacets uint32
	if err := binary.Read(bufr, binary.LittleEndian, &nfacets); err != nil {
		return nil, fmt.Errorf("error decoding STL: %v", err)
	}

	small := -math.MaxFloat64
	big := math.MaxFloat64
	min, max := Vertex3{big, big, big}, Vertex3{small, small, small}
	facets := make([]*facet, nfacets)
	for i := range facets {
		bufr.Discard(12) // discard normal
		var vertices [3]Vertex3
		for vi := range vertices {
			v, err := getVertex(bufr)
			if err != nil {
				return nil, fmt.Errorf("error decoding STL: %v", err)
			}
			vertices[vi] = v
			min.X = math.Min(min.X, v.X)
			max.X = math.Max(max.X, v.X)
			min.Y = math.Min(min.Y, v.Y)
			max.Y = math.Max(max.Y, v.Y)
			min.Z = math.Min(min.Z, v.Z)
			max.Z = math.Max(max.Z, v.Z)
		}
		facets[i] = &facet{
			vertices: vertices,
			lowZ:     math.Min(math.Min(vertices[0].Z, vertices[1].Z), vertices[2].Z),
			highZ:    math.Max(math.Max(vertices[0].Z, vertices[1].Z), vertices[2].Z),
		}
		if _, err := bufr.Discard(2); err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}
	}

	s := STL{
		facets: facets,
		Min:    min,
		Max:    max,
	}
	return &s, nil
}

func getVertex(r io.Reader) (Vertex3, error) {
	var x, z, y float32
	err := binary.Read(r, binary.LittleEndian, &x)
	if err != nil {
		return Vertex3{}, err
	}
	err = binary.Read(r, binary.LittleEndian, &y)
	if err != nil {
		return Vertex3{}, err
	}
	err = binary.Read(r, binary.LittleEndian, &z)
	if err != nil {
		return Vertex3{}, err
	}

	v := Vertex3{X: float64(x), Y: float64(y), Z: float64(z)}
	return v, nil
}
