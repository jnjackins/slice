package slice

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestImage(t *testing.T) {
	f, err := os.Open("./testdata/pikachu.stl")
	if err != nil {
		t.Fatal(err)
	}
	stl, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	var cfg = Config{
		LayerHeight: 1.0,
	}

	err = stl.Slice(nil, cfg)
	if err != nil {
		t.Error(err)
	}

	for i := range stl.Layers {
		img := stl.Layers[i].Image()
		f, err = os.Create(fmt.Sprintf("./testdata/out%d.png", i))
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, img); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
}
