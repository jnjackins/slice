package slice

import (
	"os"
	"testing"

	"sigint.ca/slice/stl"
)

func TestSlice(t *testing.T) {
	t.Log("opening stl file")
	f, err := os.Open("testdata/pikachu.stl")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("parsing stl")
	stl, err := stl.Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	var cfg = Config{
		LayerHeight: 1.0,
		LineWidth:   1.0,
	}
	layers, err := Slice(stl, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("sliced %d layers", len(layers))
}
