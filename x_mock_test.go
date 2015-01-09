package x

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// MockMultipart mocks a multipart file reading from a buffered slice.
type MockMultipart bytes.Buffer

// Read from the underlying bytes.Buffer.
func (m MockMultipart) Read(p []byte) (int, error) {
	b := bytes.Buffer(m)
	return (&b).Read(p)
}

// ReadAt is a nop.
func (m MockMultipart) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

// Seek is a nop.
func (m MockMultipart) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

// Close is a nop.
func (m MockMultipart) Close() error {
	return nil
}

// MockStore mocks an image store.
type MockStore struct{}

// Put returns a WriteCloser that writes to /dev/null and a Close method that
// is a nop.
func (m MockStore) Put(path string) (io.WriteCloser, error) {
	return struct {
		io.Writer
		io.Closer
	}{
		ioutil.Discard,
		ioutil.NopCloser(bytes.NewReader([]byte("nop"))),
	}, nil
}

// MockFileStore mocks a file store that allows you to store a file at a given
// path.
type MockFileStore struct{}

// Put returns a WriteCloser that writes to the given file path on the local
// file system.
// NOTE Will return an error if a directory in the path is not created.
func (m MockFileStore) Put(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
