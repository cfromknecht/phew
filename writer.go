package phew

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"io"

	"github.com/davecgh/go-spew/spew"
)

// A Writer is an io.WriteCloser used to encode binary data into a gzipped,
// hex-encoded binary stream.
type Writer struct {
	gw *gzip.Writer
}

// Write writes a compressed, hex-encoded form of p to the underlying io.Writer.
// The compressed bytes are not necessarily flushed until the Writer is closed.
func (w *Writer) Write(p []byte) (int, error) {
	return w.gw.Write(p)
}

// NewWriter returns a new Writer. Writes to the returned writer are compressed,
// hex-encoded and written to w.

// It is the caller's responsibility to call Close on the Writer when done.
// Writes may be buffered and not flushed until Close.
func NewWriter(w io.Writer) *Writer {
	z := new(Writer)
	z.Reset(w)
	return z
}

// Reset discards the Writer z's state and makes it equivalent to the result of
// its original state from NewWriter, but writing to w instead. This permits
// reusing a Writer rather than allocating a new one.
func (z *Writer) Reset(w io.Writer) {
	*z = Writer{gw: gzip.NewWriter(hex.NewEncoder(w))}
}

// Close closes the Writer by flushing any unwritten data to the underlying
// io.Writer and writing the hex-encoded, GZIP footer. It does not close the
// underlying io.Writer.
func (z *Writer) Close() error {
	return z.gw.Close()
}

// Flush flushes any pending compressed data to the underlying writer.

// It is useful mainly in compressed network protocols, to ensure that a remote
// reader has enough data to reconstruct a packet. Flush does not return until
// the data has been written. If the underlying writer returns an error, Flush
// returns that error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.
func (z *Writer) Flush() error {
	return z.gw.Flush()
}

// Sdump spews a varidic number of objects into a phew stream.
func Sdump(v ...interface{}) string {
	var b bytes.Buffer
	w := NewWriter(&b)
	spew.Fdump(w, v...)
	_ = w.Close()

	return string(b.Bytes())
}
