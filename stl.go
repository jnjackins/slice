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
	Layers     []*Layer
	facets     []*facet
	minZ, maxZ float64
}

type facet struct {
	normal      vertex //TODO: ignore?
	vertices    [3]vertex
	lowZ, highZ float64
}

func (f *facet) String() string {
	return fmt.Sprint(f.vertices)
}

type vertex struct {
	x, y, z float64
}

func (v vertex) String() string {
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
	facets := make([]*facet, nfacets)
	for i := range facets {
		normal, err := getVertex(f)
		if err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}
		var vertices [3]vertex
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
		facets[i] = &facet{
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

func getVertex(r io.Reader) (vertex, error) {
	var x, z, y float32
	err := binary.Read(r, binary.LittleEndian, &x)
	if err != nil {
		return vertex{}, err
	}
	err = binary.Read(r, binary.LittleEndian, &y)
	if err != nil {
		return vertex{}, err
	}
	err = binary.Read(r, binary.LittleEndian, &z)
	if err != nil {
		return vertex{}, err
	}

	v := vertex{x: float64(x), y: float64(y), z: float64(z)}
	return v, nil
}
