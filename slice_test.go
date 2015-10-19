package slice

import (
	"os"
	"testing"
)

func TestSlice(t *testing.T) {
	t.Log("opening stl file")
	f, err := os.Open("./testdata/pikachu.stl")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("parsing stl")
	stl, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	var cfg = Config{
		LayerHeight:   1.0,
		LineWidth:     1.0,
		InfillSpacing: 2.0,
		InfillAngle:   45.0,
	}

	t.Log("slicing with output=./testdata/output.gcode")
	out, err := os.Create("./testdata/output.gcode")
	if err != nil {
		t.Fatal(err)
	}
	err = stl.Slice(out, cfg)
	if err != nil {
		t.Fatal(err)
	}
	out.Close()
}
