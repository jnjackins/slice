package slice

import (
	"os"
	"testing"
)

func TestSlice(t *testing.T) {
	f, err := os.Open("./testdata/cube40_binary.stl")
	if err != nil {
		t.Fatal(err)
	}
	stl, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	var cfg = Config{
		LayerHeight: 1.0,
	}

	err = stl.Slice(os.Stdout, cfg)
	if err != nil {
		t.Error(err)
	}
}
