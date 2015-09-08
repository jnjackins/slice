// Package slice provides types and functions for compiling STL format 3D models
// into G-code to be used for 3D printing
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

// Slice coordinates parallel slicing (1 goroutine per layer)
func (s *STL) Slice(w io.Writer, cfg Config) error {
	var wg sync.WaitGroup
	nLayers := int((s.maxZ-s.minZ)/cfg.LayerHeight) + 1
	dprintf("sliced %d layers", nLayers)
	s.layers = make([]*Layer, nLayers)
	for i := range s.layers {
		wg.Add(1)
		//TODO: go func
		func(i int, z float64) {
			s.layers[i] = s.mkLayer(z)
			wg.Done()
		}(i, float64(i)*cfg.LayerHeight)
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
