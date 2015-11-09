package slice

import (
	"bufio"
	"bytes"
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

	var getVertex func(r *bufio.Reader) (Vertex3, error)
	var facets []*facet
	small := math.Inf(-1)
	big := math.Inf(+1)
	min, max := Vertex3{big, big, big}, Vertex3{small, small, small}

	getFacet := func() (*facet, error) {
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
		return &facet{
			vertices: vertices,
			lowZ:     math.Min(math.Min(vertices[0].Z, vertices[1].Z), vertices[2].Z),
			highZ:    math.Max(math.Max(vertices[0].Z, vertices[1].Z), vertices[2].Z),
		}, nil
	}

	if string(top) == "solid " {
		getVertex = getVertexAscii

		// discard header
		bufr.ReadLine()

		facets = make([]*facet, 0)
		for {
			bufr.ReadLine() // discard normal
			bufr.ReadLine() // discard "outer loop"

			f, err := getFacet()
			if err != nil {
				return nil, err
			}
			facets = append(facets, f)

			bufr.ReadLine() // discard "endloop"
			bufr.ReadLine() // discard "endfacet"

			nextWord, err := bufr.Peek(8)
			if err != nil {
				return nil, fmt.Errorf("error decoding STL: %v", err)
			}
			if string(nextWord) == "endsolid" {
				break
			}
		}
	} else {
		getVertex = getVertexBinary

		// discard header
		if _, err := bufr.Discard(80); err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}

		var nfacets uint32
		if err := binary.Read(bufr, binary.LittleEndian, &nfacets); err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}

		facets = make([]*facet, nfacets)
		for i := range facets {
			bufr.Discard(12) // discard normal

			f, err := getFacet()
			if err != nil {
				return nil, err
			}
			facets[i] = f

			if _, err := bufr.Discard(2); err != nil {
				return nil, fmt.Errorf("error decoding STL: %v", err)
			}
		}
	}

	s := STL{
		facets: facets,
		Min:    min,
		Max:    max,
	}
	return &s, nil
}

func getVertexAscii(r *bufio.Reader) (Vertex3, error) {
	var x, z, y float32

	// sometimes ASCII STLs are indented, sometimes they aren't. strip leading whitespace
	// if it exists.
	s, _, err := r.ReadLine()
	if err != nil {
		return Vertex3{}, err
	}
	bytes.TrimSpace(s)

	if _, err := fmt.Sscanf(string(s), "vertex %f %f %f\n", &x, &y, &z); err != nil {
		return Vertex3{}, err
	}
	v := Vertex3{X: float64(x), Y: float64(y), Z: float64(z)}
	return v, nil
}

func getVertexBinary(r *bufio.Reader) (Vertex3, error) {
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
