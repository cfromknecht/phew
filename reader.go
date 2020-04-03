package phew

import (
	"compress/gzip"
	"encoding/hex"
	"io"
)

// A Reader is an io.Reader that can be read to retrieve hex-decoded,
// uncompressed data from a phew-format compressed file.

// In general, a phew file can be a concatenation of phew files, each with its
// own header. Reads from the Reader return the concatenation of the
// uncompressed data of each.

// Hzip files store a length and checksum of the uncompressed data. The Reader
// will return an ErrChecksum when Read reaches the end of the uncompressed data
// if it does not have the expected length or checksum. Clients should treat
// data returned by Read as tentative until they receive the io.EOF marking the
// end of the data.
type Reader struct {
	gr *gzip.Reader
}

// NewReader creates a new Reader reading the given reader. If r does not also
// implement io.ByteReader, the decompressor may read more data than necessary
// from r.
//
// It is the caller's responsibility to call Close on the Reader when done.
func NewReader(r io.Reader) (*Reader, error) {
	z := new(Reader)
	if err := z.Reset(r); err != nil {
		return nil, err
	}

	return z, nil
}

// Reset discards the Reader z's state and makes it equivalent to the result of
// its original state from NewReader, but reading from r instead. This permits
// reusing a Reader rather than allocating a new one.
func (z *Reader) Reset(r io.Reader) error {
	gr, err := gzip.NewReader(hex.NewDecoder(r))
	if err != nil {
		return err
	}

	*z = Reader{gr: gr}

	return nil
}

// Read implements io.Reader, reading uncompressed bytes from its underlying
// Reader.
func (z *Reader) Read(p []byte) (int, error) {
	return z.gr.Read(p)
}
