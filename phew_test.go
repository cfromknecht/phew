package phew_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/cfromknecht/phew"
	"github.com/davecgh/go-spew/spew"
)

func TestEncodeDecode(t *testing.T) {
	t.Run("no flush", func(t *testing.T) {
		testEncodeDecode(t, false)
	})
	t.Run("flush", func(t *testing.T) {
		testEncodeDecode(t, true)
	})
}

func testEncodeDecode(t *testing.T, flush bool) {
	testStr := []byte("hello world")

	var b bytes.Buffer
	w := phew.NewWriter(&b)

	_, err := w.Write(testStr)
	if err != nil {
		t.Fatalf("unable to write: %v", err)
	}

	if flush {
		err = w.Flush()
		if err != nil {
			t.Fatalf("unable to flush: %v", err)
		}
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("unable to close: %v", err)
	}

	r, err := phew.NewReader(bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatalf("unable to create reader: %v", err)
	}

	decStr, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("unable to read: %v", err)
	}

	if !bytes.Equal(testStr, decStr) {
		t.Fatalf("mismatch, expected: %v, got: %v", testStr, decStr)
	}
}

func TestDecodeInvalidHeader(t *testing.T) {
	_, err := phew.NewReader(bytes.NewReader([]byte{}))
	if err == nil {
		t.Fatalf("reader unexpectedly initialized")
	}
}

func TestSdump(t *testing.T) {
	type VeryLargeThing struct {
		Number      uint32
		LotsOfBytes []byte
	}

	ting := VeryLargeThing{42, []byte("big data")}

	plaintext := spew.Sdump(ting)

	var b bytes.Buffer
	w := phew.NewWriter(&b)

	_, err := w.Write([]byte(plaintext))
	if err != nil {
		t.Fatalf("unable to write: %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("unable to close: %v", err)
	}

	enc := b.Bytes()
	expEnc := []byte(phew.Sdump(ting))

	if !bytes.Equal(enc, expEnc) {
		t.Fatalf("mismatch, expected: %v, got: %v", expEnc, enc)
	}
}
