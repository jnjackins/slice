// TODO: allow the client to slice, infill, e.g. independently

// package slice provides types and functions for slicing and compiling STL format 3D models
// into G-code to be used for 3D printing.
package slice

import (
	"fmt"
	"os"
	"sync"

	"sigint.ca/slice/stl"
)

var debug bool

// A Config variable specifies a slicing configuration.
type Config struct {
	DebugMode bool

	LayerHeight float64
	LineWidth   float64

	InfillSpacing float64
	InfillAngle   float64 // in degrees
}

func dprintf(format string, args ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, "[ "+format+" ]\n", args...)
	}
}

func wprintf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "WARNING: "+format+"\n", args...)
}

// Slice slices and stl.Solid into layers.
func Slice(s *stl.Solid, cfg Config) ([]*Layer, error) {
	debug = cfg.DebugMode

	min, max := s.Bounds()
	nLayers := int(0.5 + (max.Z-min.Z)/cfg.LayerHeight)
	layers := make([]*Layer, nLayers)
	h := cfg.LayerHeight

	// slice in parallel if not in debug mode
	if debug {
		for i := range layers {
			layers[i] = sliceLayer(i, min.Z+0.01+float64(i)*h, s, cfg)
			//layers[i].genInfill(cfg)
		}
	} else {
		var wg sync.WaitGroup
		for i := range layers {
			wg.Add(1)
			go func(i int, z float64) {
				layers[i] = sliceLayer(i, z, s, cfg)
				//layers[i].genInfill(cfg)
				wg.Done()
			}(i, min.Z+0.01+float64(i)*h)
		}
		wg.Wait()
	}

	dprintf("sliced %d layers", nLayers)
	return layers, nil
}
