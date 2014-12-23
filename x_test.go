package x

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"
)

// TestMock tests if all the required methods are implemented for mocking.
func TestMock(t *testing.T) {
	b := MockMultipart(*bytes.NewBuffer([]byte("test")))
	_ = NewJob(b, Size{}, PNG, MockStore{})
}

// TestProcess runs image processing on a test.jpeg and tests the computed SHA1
// for the image.
func TestProcess(t *testing.T) {
	in, err := os.Open("test.jpeg")
	if err != nil {
		t.Fatalf("error opening test image : %s", err)
	}
	j := NewJob(in, Size{Square, 80}, WebP, MockFileStore{})
	if err := j.Process(); err != nil {
		t.Fatalf("error while processing image : %s", err)
	}
	sha := fmt.Sprintf("%x", j.input.Sum(nil))
	if x := "54c23af026a57bad3d146cee288cf8f1971ca2c7"; sha != x {
		t.Errorf("computed sha got %s, expected %s", sha, x)
	}
}

// sizeTests values for table driven tests of the size method.
var sizeTests = []struct {
	// size is input size constraint.
	size Size
	// input and expected ouput width and height.
	inW, inH, outW, outH int
	// error that is expected.
	err error
}{
	// Square sizing
	{Size{Square, 40}, 80, 160, 40, 40, nil},
	{Size{Square, 40}, 160, 80, 40, 40, nil},
	{Size{Square, 40}, 20, 80, 20, 40, nil},
	{Size{Square, 40}, 20, 40, 20, 40, nil},
	// Max sizing
	{Size{Max, 40}, 80, 160, 20, 40, nil},
	{Size{Max, 40}, 160, 80, 40, 20, nil},
	{Size{Max, 40}, 20, 80, 10, 40, nil},
	{Size{Max, 40}, 20, 40, 20, 40, nil},
	// MaxHeight sizing
	{Size{MaxHeight, 40}, 80, 160, 20, 40, nil},
	{Size{MaxHeight, 40}, 160, 80, 80, 40, nil},
	{Size{MaxHeight, 40}, 80, 20, 80, 20, nil},
	{Size{MaxHeight, 40}, 40, 20, 40, 20, nil},
	// MaxWidth sizing
	{Size{MaxWidth, 40}, 80, 160, 40, 80, nil},
	{Size{MaxWidth, 40}, 160, 80, 40, 20, nil},
	{Size{MaxWidth, 40}, 20, 80, 20, 80, nil},
	{Size{MaxWidth, 40}, 20, 40, 20, 40, nil},
	// Invalid input
	{Size{SizeType("invalidsizetype"), 40}, 20, 40, 20, 40,
		ErrInvalidSizeType},
}

// TestSize is a table driven test for the various sizing methods.
func TestSize(t *testing.T) {
	for _, tt := range sizeTests {
		pre := fmt.Sprintf("%d x %d input under %s %d constraint : ",
			tt.inW, tt.inH, tt.size.Type, tt.size.Value)
		m := image.NewRGBA(image.Rect(0, 0, tt.inW, tt.inH))
		out, err := size(m, tt.size)
		if err != tt.err {
			t.Errorf(pre+"error %s, expected %s", err.Error(),
				tt.err.Error())
		}
		if tt.err == nil {
			if out.Bounds().Dx() != tt.outW {
				t.Errorf(pre+"width %d, expected %d",
					out.Bounds().Max.X, tt.outW)
			}
			if out.Bounds().Dy() != tt.outH {
				t.Errorf(pre+"height %d, expected %d",
					out.Bounds().Max.Y, tt.outH)
			}
		}
	}
}
