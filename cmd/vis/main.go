package main

import (
	"image"
	"image/draw"
	"log"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"sigint.ca/slice"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("vis: ")
	f, err := os.Open("../../testdata/pikachu.stl")
	if err != nil {
		log.Fatal(err)
	}
	stl, err := slice.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	var cfg = slice.Config{
		LayerHeight: 0.4,
	}

	err = stl.Slice(nil, cfg)
	if err != nil {
		log.Fatal(err)
	}

	var layerImages = make([]*image.RGBA, len(stl.Layers))
	for i := range stl.Layers {
		layerImages[i] = stl.Layers[i].Draw()
	}
	layer := 0

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		r := image.Rect(int(stl.Min.X), int(stl.Min.Y), int(stl.Max.X)+1, int(stl.Max.Y)+1)
		winSize := r.Size()
		var b screen.Buffer
		defer func() {
			if b != nil {
				b.Release()
			}
		}()

		var sz size.Event
		var lastClick mouse.Event

		for e := range w.Events() {
			switch e := e.(type) {
			default:

			case mouse.Event:
				if e.Button == mouse.ButtonLeft {
					if e.Y > lastClick.Y && layer < len(layerImages)-1 {
						layer++
					} else if e.Y < lastClick.Y && layer > 0 {
						layer--
					} else {
						break
					}
					lastClick = e

					src, dst := layerImages[layer], b.RGBA()
					draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case paint.Event:
				w.Upload(image.Point{}, b, b.Bounds(), w)
				w.EndPaint(e)

			case screen.UploadedEvent:
				// No-op.

			case size.Event:
				sz = e
				if b != nil {
					b.Release()
				}
				winSize = image.Point{sz.WidthPx, sz.HeightPx}
				b, err = s.NewBuffer(winSize)
				if err != nil {
					log.Fatal(err)
				}
				src, dst := layerImages[layer], b.RGBA()
				draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

			case error:
				log.Printf("error: %v", e)
			}
		}
	})
}
