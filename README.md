Slice
=====

Slice is a STL to G-code compiler for 3D printers.

[Package documentation](https://godoc.org/sigint.ca/slice)

## Goals:
* **Fast** - never wait more than 1 second
* **Robust** - never crash
* **Multi-platform** - run on any platform supported by [Go](https://golang.org/doc/install#requirements),
with minimal dependencies

## What works:
* Perimeter slicing
* Basic linear infill
* Sliced layer previews

## What doesn't work:
* Printing with generated G-code

## Try it out
```
go get -t -u sigint.ca/slice/cmd/preview
cd $GOPATH/src/sigint.ca/slice
go build sigint.ca/slice/cmd/preview
./preview testdata/pikachu.stl
```

Click and drag up and down to scroll through layers.
