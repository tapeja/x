// Package x is under construction.
package x

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"hash"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"

	"github.com/chai2010/gopkg/image/webp"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// Size provides a dimension constraint on an image.
type Size struct {
	// Type is the contraint type.
	Type SizeType
	// Value is the constraint parameter, usually in pixels.
	Value int
}

// SizeType provides enum values for the various contraint types.
type SizeType string

const (
	// Square mode creates a square image with the dimensions of the passed
	// parameter. First crops the larger dimension equally from both sides
	// to match the smaller dimension then shrinks the square down to the
	// passed parameter. Image isn't resized if both dimensions are smaller
	// than the passed parameter.
	Square SizeType = "square"
	// Max mode determines the larger dimension of the image and shrinks it
	// to match the passed parameter retaining the original aspect ratio.
	// Image isn't resized if both dimensions are smaller than the passed
	// parameter.
	Max = "max"
	// MaxWidth is just like Max but fixed on the width. Image isn't
	// resized if the fixed dimension is smaller than the passed parameter.
	MaxWidth = "max_width"
	// MaxHeight is just like Max but fixed on the height. Image isn't
	// resized if the fixed dimension is smaller than the passed parameter.
	MaxHeight = "max_height"
)

// Format provides a file format to output.
type Format string

const (
	JPG  Format = "jpg"
	PNG  Format = "png"
	WebP Format = "webp"
)

// Store stores images for later retrieval.
type Store interface {
	Put(path string) (out io.WriteCloser, err error)
}

// Job holds the configuration and inputs for an image processing job.
//
// Steps to process an image :
//
//	 1. Calculate the SHA1 of the original image.
//	 2. Decode the image from the original format.
//	 3. Resize the image based on constraints.
//	 4. Encode into the output format.
//	 5. Write to the provided store.
//
type Job struct {
	input  *input
	size   Size
	format Format
	store  Store
}

// NewJob creates a new image processing job.
func NewJob(r io.ReadCloser, size Size, format Format, store Store) *Job {
	return &Job{
		input:  newInput(r),
		size:   size,
		format: format,
		store:  store,
	}
}

// ErrInvalidFormat is returned when attempting to encode an image into an
// unsupported file format.
var ErrInvalidFormat error = errors.New("invalid image format")

// Process initiates an image processing job.
func (j *Job) Process() (err error) {
	// Close and override the returned err only if not set already.
	closeErr := func(c io.Closer) {
		if c != nil {
			if cerr := c.Close(); err == nil {
				err = cerr
			}
		}
	}

	// Process the input
	defer closeErr(j.input)
	m, f, err := image.Decode(j.input)
	if err != nil {
		return err
	}
	log.Printf("Decoded image with format %s and dimensions %d x %d\n",
		f, m.Bounds().Dx(), m.Bounds().Dy())
	// Some codecs don't seem to read full file flush to compute SHA.
	io.Copy(ioutil.Discard, j.input)
	sha := fmt.Sprintf("%x", j.input.Sum(nil))
	log.Printf("Computed original image SHA1 as %s\n", sha)

	// Size the image
	processed, err := size(m, j.size)
	if err != nil {
		return err
	}
	log.Printf("Sized image to %d x %d\n", processed.Bounds().Dx(),
		processed.Bounds().Dy())

	// Encode and write the file to the store
	// TODO Later change to path with directories, to avoid file system
	// limits.
	path := fmt.Sprintf("%s-%d-%s.%s", j.size.Type, j.size.Value, sha,
		j.format)
	out, err := j.store.Put(path)
	defer closeErr(out)
	if err != nil {
		return err
	}
	log.Printf("Writing image to %s", path)
	switch j.format {
	case JPG:
		return jpeg.Encode(out, processed, &jpeg.Options{100})
	case PNG:
		return png.Encode(out, processed)
	case WebP:
		return webp.Encode(out, processed, &webp.Options{true, 100})
	default:
		return ErrInvalidFormat
	}
}

// ErrInvalidSizeType is thrown when the processing job is provided an invalid
// SizeType.
var ErrInvalidSizeType error = errors.New("invalid size constraint")

// size applies the passed size contraint on the passed decoded image, returns
// sized image and any errors encountered during sizing.
func size(image image.Image, size Size) (image.Image, error) {
	switch size.Type {
	case Square:
		return sizeSquare(image, size.Value)
	case Max:
		return sizeMax(image, size.Value)
	case MaxWidth:
		return sizeMaxWidth(image, size.Value)
	case MaxHeight:
		return sizeMaxHeight(image, size.Value)
	default:
		return nil, ErrInvalidSizeType
	}
}

// sizeSquare applies the Square size constraint on the image, first cropping
// the larger dimensions to match the smaller dimension then resizing to a
// square of the passed dimension.
func sizeSquare(image image.Image, dim int) (image.Image, error) {
	w, h := image.Bounds().Dx(), image.Bounds().Dy()
	var crop, nw, nh int
	if w <= dim && h <= dim {
		// No resize, return untouched
		return image, nil
	} else if w > h {
		crop = h
		nw = dim
	} else {
		crop = w
		nh = dim
	}
	if crop < dim {
		crop = dim
	}
	// Crop to make both dimensions equal
	m, err := cutter.Crop(image, cutter.Config{
		Width:  crop,
		Height: crop,
		Mode:   cutter.Centered,
	})
	if err != nil {
		return nil, err
	}
	// Resize the image down to the parameter value
	return resize.Resize(uint(nw), uint(nh), m, resize.Lanczos3), nil
}

// sizeMax applies the Max size constraint, returning the resized image where
// the maximum dimensions is the passed the
func sizeMax(image image.Image, dim int) (image.Image, error) {
	if image.Bounds().Dx() <= dim && image.Bounds().Dy() <= dim {
		return image, nil
	}
	return resize.Thumbnail(uint(dim), uint(dim), image,
		resize.Lanczos3), nil
}

// sizeMaxWidth applies the MaxWidth constraint, returning the resized image
// with the passed maximum width.
func sizeMaxWidth(image image.Image, dim int) (image.Image, error) {
	if image.Bounds().Dx() <= dim {
		return image, nil
	}
	return resize.Resize(uint(dim), 0, image,
		resize.Lanczos3), nil
}

// sizeMaxHeight applies the MaxHeight constraint, returning the resized image
// with the passed maximum height.
func sizeMaxHeight(image image.Image, dim int) (image.Image, error) {
	if image.Bounds().Dy() <= dim {
		return image, nil
	}
	return resize.Resize(0, uint(dim), image,
		resize.Lanczos3), nil
}

// input represents the state while processing an input image.
type input struct {
	io.ReadCloser
	hash.Hash
}

// newInput creates a new input processing state based on the passed ReadCloser.
func newInput(r io.ReadCloser) *input {
	return &input{r, sha1.New()}
}

// Read method allows us to pass the input directly to the image package's
// Decode function, so that the image input is read only once for all
// processing needs such as decoding and computing the SHA. May return errors
// from hashing.
func (i *input) Read(p []byte) (n int, err error) {
	n, err = i.ReadCloser.Read(p)
	// Make sure to trim the buffer to the read bytes.
	if wn, werr := i.Write(p[:n]); werr != nil {
		return wn, werr
	}
	return
}
