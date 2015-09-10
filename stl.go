package slice

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type STL struct {
	Layers   []*Layer
	Min, Max Vertex

	facets []*facet
}

type facet struct {
	normal      Vertex //TODO: ignore?
	vertices    [3]Vertex
	lowZ, highZ float64
}

func (f *facet) String() string {
	return fmt.Sprint(f.vertices)
}

type Vertex struct {
	X, Y, Z float64
}

func (v Vertex) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f,%0.3f)", v.X, v.Y, v.Z)
}

func Parse(f *os.File) (*STL, error) {
	// test for ascii stl format
	//TODO: just read 6 bytes
	r := bufio.NewReader(f)
	if line, err := r.ReadString('\n'); err == nil {
		if strings.HasPrefix(line, "solid ") {
			return nil, fmt.Errorf("ascii format STL not supported")
		}
	}

	// discard text header
	if _, err := f.Seek(80, 0); err != nil {
		return nil, fmt.Errorf("error decoding STL: %v", err)
	}

	var nfacets uint32
	if err := binary.Read(f, binary.LittleEndian, &nfacets); err != nil {
		return nil, fmt.Errorf("error decoding STL: %v", err)
	}

	small := -1 * math.MaxFloat64
	big := math.MaxFloat64
	min, max := Vertex{big, big, big}, Vertex{small, small, small}
	facets := make([]*facet, nfacets)
	for i := range facets {
		normal, err := getVertex(f)
		if err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}
		var vertices [3]Vertex
		for vi := range vertices {
			v, err := getVertex(f)
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
		if _, err := f.Seek(2, 1); err != nil {
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

func getVertex(r io.Reader) (Vertex, error) {
	var x, z, y float32
	err := binary.Read(r, binary.LittleEndian, &x)
	if err != nil {
		return Vertex{}, err
	}
	err = binary.Read(r, binary.LittleEndian, &y)
	if err != nil {
		return Vertex{}, err
	}
	err = binary.Read(r, binary.LittleEndian, &z)
	if err != nil {
		return Vertex{}, err
	}

	v := Vertex{X: float64(x), Y: float64(y), Z: float64(z)}
	return v, nil
}
