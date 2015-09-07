// TODO: adjust vertices so minZ is at 0.0

package slice

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
)

const debug = false

var cfg = struct {
	layerHeight float64
}{
	0.2,
}

type Vertex struct {
	x, y, z float64
}

func (v Vertex) String() string {
	return fmt.Sprintf("(%0.3f,%0.3f,%0.3f)", v.x, v.y, v.z)
}

type Segment struct {
	end1, end2 Vertex
}

var ZS = Segment{}

func (s Segment) String() string {
	return fmt.Sprintf("%v-%v", s.end1, s.end2)
}

type Facet struct {
	normal      Vertex //TODO: ignore?
	vertices    [3]Vertex
	lowZ, highZ float64
}

func (f *Facet) String() string {
	return fmt.Sprintf("%v", f.vertices)
}

type Layer struct {
	facets   []*Facet
	segments []Segment
}

func (l *Layer) Gcode() []byte {
	if len(l.segments) == 0 {
		return []byte{}
	}
	buf := new(bytes.Buffer)
	// first perimeters
	s := l.segments[0]
	fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.3f E%.3f\n", s.end1.x, s.end1.y, s.end1.z, 0.0)
	for _, s := range l.segments {
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.5f E%.5f\n", s.end2.x, s.end2.y, s.end2.z, 0.0)
	}
	return buf.Bytes()
}

type STL struct {
	facets     []*Facet
	layers     []*Layer
	minZ, maxZ float64
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

func (s *STL) String() string {
	var out string
	for _, f := range s.facets {
		out += fmt.Sprintln(f)
	}
	return out
}

func (s *STL) Slice(w io.Writer) error {
	var wg sync.WaitGroup
	nLayers := int((s.maxZ-s.minZ)/cfg.layerHeight) + 1
	dprintf("sliced %d layers", nLayers)
	s.layers = make([]*Layer, nLayers)
	for i := range s.layers {
		wg.Add(1)
		//TODO: go func
		func(i int, z float64) {
			s.layers[i] = s.mkLayer(z)
			wg.Done()
		}(i, float64(i)*cfg.layerHeight)
	}
	wg.Wait()
	for _, l := range s.layers {
		_, err := w.Write(l.Gcode())
		if err != nil {
			return err
		}
	}
	return nil
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

	segments := make([]Segment, 0, len(facets))
	for _, f := range facets {
		s := mkSegment(f, z)
		if s != ZS {
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

//TODO: case where segment is one of the edges of the triangle
func mkSegment(f *Facet, z float64) Segment {
	var ends [3]Vertex
	var i int
	v := f.vertices
	// two of these cases will usually be true
	if (v[0].z > z && v[1].z < z) || (v[0].z < z && v[1].z > z) {
		x1, x2 := v[0].x, v[1].x
		y1, y2 := v[0].y, v[1].y
		z1, z2 := v[0].z, v[1].z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = Vertex{x, y, z}
		i++
	}
	if (v[0].z > z && v[2].z < z) || (v[0].z < z && v[2].z > z) {
		x1, x2 := v[0].x, v[2].x
		y1, y2 := v[0].y, v[2].y
		z1, z2 := v[0].z, v[2].z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = Vertex{x, y, z}
		i++
	}
	if (v[1].z > z && v[2].z < z) || (v[1].z < z && v[2].z > z) {
		x1, x2 := v[1].x, v[2].x
		y1, y2 := v[1].y, v[2].y
		z1, z2 := v[1].z, v[2].z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = Vertex{x, y, z}
		i++
	}

	// otherwise, a segment of the facet or the entire facet should coincide with
	// the slice plane
	if i == 0 {
		//TODO
		return Segment{}
	} else if i != 2 {
		log.Printf("warning: found %d intersections when finding segment at z=%f", i, z)
		return Segment{}
	}

	return Segment{ends[0], ends[1]}
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

func dprintf(format string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, "[ "+format+" ]\n", args...)
	}
}
