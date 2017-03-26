package stl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"

	"io"
)

type binaryReader struct {
	r *bufio.Reader
}

func (p binaryReader) readFacets() ([]Facet, error) {
	// discard header
	if _, err := p.r.Discard(80); err != nil {
		return nil, fmt.Errorf("error decoding STL: %v", err)
	}

	var nfacets uint32
	if err := binary.Read(p.r, binary.LittleEndian, &nfacets); err != nil {
		return nil, fmt.Errorf("error decoding STL: %v", err)
	}

	facets := make([]Facet, nfacets)
	var vbuf = make([]byte, 12) // reusable buffer for reading vertices
	for i := range facets {
		p.r.Discard(12) // discard normal

		var f Facet
		if err := p.readFacet(&f, vbuf); err != nil {
			return nil, err
		}
		facets[i] = f

		if _, err := p.r.Discard(2); err != nil {
			return nil, fmt.Errorf("error decoding STL: %v", err)
		}
	}

	return facets, nil
}

func (p binaryReader) readFacet(f *Facet, vbuf []byte) error {
	for i := range f.Vertices {
		_, err := io.ReadFull(p.r, vbuf)
		if err != nil {
			return err
		}
		f.Vertices[i].X = float64(math.Float32frombits(binary.LittleEndian.Uint32(vbuf[0:4])))
		f.Vertices[i].Y = float64(math.Float32frombits(binary.LittleEndian.Uint32(vbuf[4:8])))
		f.Vertices[i].Z = float64(math.Float32frombits(binary.LittleEndian.Uint32(vbuf[8:12])))
	}
	f.calculateNormal()
	return nil
}
