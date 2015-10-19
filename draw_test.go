package slice

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestDraw(t *testing.T) {
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

	t.Log("slicing with output=nil")
	err = stl.Slice(nil, cfg)
	if err != nil {
		t.Error(err)
	}

	r := stl.Layers[0].Bounds()

	t.Log("drawing layers to png")
	for i := range stl.Layers {
		img := image.NewRGBA(r)
		stl.Layers[i].Draw(img)
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
