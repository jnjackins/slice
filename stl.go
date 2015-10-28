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
	normal      Vertex3 //TODO: ignore?
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

	small := -1 * math.MaxFloat64
	big := math.MaxFloat64
	min, max := Vertex3{big, big, big}, Vertex3{small, small, small}
	facets := make([]*facet, nfacets)
	for i := range facets {
		normal, err := getVertex(bufr)
		if err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}
		var vertices [3]Vertex3
		for vi := range vertices {
			v, err := getVertex(bufr)
			if err != nil {
				return nil, fmt.Errorf("error decoding STL: %v", err)
			}
			vertices[vi] = v

			if v.X < min.X {
				min.X = v.X
			}
			if v.X > max.X {
				max.X = v.X
			}
			if v.Y < min.Y {
				min.Y = v.Y
			}
			if v.Y > max.Y {
				max.Y = v.Y
			}
			if v.Z < min.Z {
				min.Z = v.Z
			}
			if v.Z > max.Z {
				max.Z = v.Z
			}
		}
		facets[i] = &facet{
			normal:   normal,
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
