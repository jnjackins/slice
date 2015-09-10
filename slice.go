// TODO: split 1 layer into sublayers, so that they can have different infills etc.

// Package slice provides types and functions for slicing and compiling STL format 3D models
// into G-code to be used for 3D printing.
package slice // import "sigint.ca/slice"

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const debug = false

type Config struct {
	LayerHeight float64
}

func dprintf(format string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, "[ "+format+" ]\n", args...)
	}
}

// Slice divides parsed STL data into layers and compiles G-code
// for each layer. The resulting G-code is written to w if it is not nil.
// After running slice, the resulting layers can be accessed (and compiled
// individually into G-code) by accessing the STL's Layers variable.
func (s *STL) Slice(w io.Writer, cfg Config) error {
	var wg sync.WaitGroup
	nLayers := int((s.Max.Z-s.Min.Z)/cfg.LayerHeight) + 1
	dprintf("sliced %d layers", nLayers)
	s.Layers = make([]*Layer, nLayers)
	h := cfg.LayerHeight
	for i := range s.Layers {
		wg.Add(1)
		//TODO: go func
		func(i int, z float64) {
			s.Layers[i] = s.sliceLayer(z)
			wg.Done()
		}(i, 0.001+float64(i)*h)
	}
	wg.Wait()
	if w == nil {
		return nil
	}
	for _, l := range s.Layers {
		_, err := w.Write(l.Gcode())
		if err != nil {
			return err
		}
	}
	return nil
}
