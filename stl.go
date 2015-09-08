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
	facets     []*Facet
	layers     []*Layer
	minZ, maxZ float64
}

type Facet struct {
	normal      Vertex //TODO: ignore?
	vertices    [3]Vertex
	lowZ, highZ float64
}

func (f *Facet) String() string {
	return fmt.Sprint(f.vertices)
}

type Vertex struct {
	x, y, z float64
}

func (v Vertex) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f,%0.3f)", v.x, v.y, v.z)
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

	var maxZ, minZ float64 = 0, math.MaxFloat64
	facets := make([]*Facet, nfacets)
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
			if v.z < minZ {
				minZ = v.z
			}
			if v.z > maxZ {
				maxZ = v.z
			}
		}
		facets[i] = &Facet{
			normal:   normal,
			vertices: vertices,
			lowZ:     math.Min(math.Min(vertices[0].z, vertices[1].z), vertices[2].z),
			highZ:    math.Max(math.Max(vertices[0].z, vertices[1].z), vertices[2].z),
		}
		if _, err := f.Seek(2, 1); err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}
	}

	s := STL{
		facets: facets,
		minZ:   minZ,
		maxZ:   maxZ,
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

	v := Vertex{x: float64(x), y: float64(y), z: float64(z)}
	return v, nil
}
