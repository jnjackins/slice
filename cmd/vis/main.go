package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"sigint.ca/slice"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	debug = flag.Bool("d", false, "debug mode")
	prof  = flag.String("prof", "", "`path` to output CPU profiling information")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("vis: ")
	flag.Usage = func() {
		log.Print("Usage: vis [options] file")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *prof != "" {
		log.Print("profile mode: will write out CPU profile after 5 seconds")
		f, err := os.Create(*prof)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			time.Sleep(5 * time.Second)
			pprof.StopCPUProfile()
			log.Print("done writing CPU profile")
		}()
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	stl, err := slice.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	if err := sliceSTL(stl); err != nil {
		log.Fatal(err)
	}
	imgs := drawLayers(stl)

	log.Print("Launching UI...")
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

		var layer int
		redraw := func() {
			draw.Draw(b.RGBA(), b.RGBA().Bounds(), imgs[layer], imgs[layer].Bounds().Min, draw.Src)
			drawLayerNumber(b.RGBA(), layer)
		}

		for e := range w.Events() {
			switch e := e.(type) {
			default:

			case mouse.Event:
				if e.Button == mouse.ButtonLeft {
					if e.Y > lastClick.Y && layer < len(imgs)-1 {
						layer++
						redraw()
					} else if e.Y < lastClick.Y && layer > 0 {
						layer--
						redraw()
					}
					lastClick = e
				}

			case key.Event:
				log.Printf("key: %v", e)
				if e.Code == key.CodeEscape {
					log.Print("quitting")
					return
				}

			case paint.Event:
				w.Upload(image.Point{}, b, b.Bounds(), w)
				w.Publish()

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
				redraw()

			case error:
				log.Printf("error: %v", e)
			}
		}
	})
}

func sliceSTL(stl *slice.STL) error {
	log.Print("slicing...")
	t := time.Now()

	var cfg = slice.Config{
		DebugMode:     *debug,
		LayerHeight:   0.4,
		LineWidth:     0.2,
		InfillSpacing: 1.0,
		InfillAngle:   45.0,
	}

	if err := stl.Slice(nil, cfg); err != nil {
		return err
	}

	log.Printf("slicing took %v", time.Now().Sub(t))
	return nil
}

func drawLayers(stl *slice.STL) []*image.RGBA {
	log.Print("drawing layers...")
	if len(stl.Layers) == 0 {
		log.Fatal("no layers to draw")
	}
	t := time.Now()

	imgs := make([]*image.RGBA, len(stl.Layers))
	r := stl.Layers[0].Bounds()
	wg := sync.WaitGroup{}
	for i, layer := range stl.Layers {
		wg.Add(1)
		go func(i int, layer *slice.Layer) {
			imgs[i] = image.NewRGBA(r)
			draw.Draw(imgs[i], r, image.White, r.Min, draw.Src)
			layer.Draw(imgs[i])
			wg.Done()
		}(i, layer)
	}
	wg.Wait()

	log.Printf("drawing took %v", time.Now().Sub(t))
	return imgs
}

func drawLayerNumber(dst draw.Image, n int) {
	d := font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: basicfont.Face7x13,
		Dot:  fixed.P(2, 13),
	}
	d.DrawString(fmt.Sprintf("Layer %03d", n))
}
