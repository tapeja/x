package main

// Size provides a dimension constraint on an image.
type Size string

const (
	Square    Size = "square"
	Max       Size = "max"
	MaxWidth  Size = "max_width"
	MaxHeight Size = "max_height"
)

// Format provides a file format to output.
type Format string

const (
	JPG  Format = "jpg"
	PNG  Format = "png"
	WebP Format = "webp"
)

func main() {

}
