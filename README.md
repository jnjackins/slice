Slice
=====

Slice is a set of tools and Go packages for compiling STL files into toolpaths for 3D printers.

[Package documentation](https://godoc.org/sigint.ca/slice)

## What works:
* Perimeter slicing
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
