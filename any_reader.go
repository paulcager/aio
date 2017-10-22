// Don't use standard file detection software / libmagic as it requires >= 128 bytes to be read.
// https://en.wikipedia.org/wiki/List_of_file_signatures

// ZIP files / tar files - return concat of all contained files.

package io

import (
	"compress/gzip"
	"errors"
	"io"
)

var NotSupported = errors.New("The file type is not yet supported")

const (
	compressMagic = "\x1f\x9d"
	gzipMagic     = "\x1f\x8b"
	lzipMagic     = "LZIP"
	bzip2Magic    = "BZh"
	xzMagic       = "\x37\x7A\xBC\xAF\x27\x1C"
)

type anyReader struct {
	r       io.Reader
	decided bool
}

func (r *anyReader) Read(b []byte) (n int, err error) {
	if !r.decided {
		err = r.decide()
		if err != nil {
			return 0, err
		}
	}

	return r.r.Read(b)
}

//	p.peak = p.peak[n:]
func NewAnyReader(r io.Reader) *anyReader {
	return &anyReader{r: r}
}

func (r *anyReader) decide() error {
	var err error
	if r.decided {
		return nil
	}

	peaker := NewPeakableReader(r.r)
	r.r = peaker
	b := make([]byte, 512)
	r.decided = true

	if n, _ := peaker.PeakAtLeast(b, 2); n == 2 && string(b[:2]) == string([]byte{0x1f, 0x9d}) {
		// "compress" format. https://en.wikipedia.org/wiki/Lempel-Ziv-Welch
		return NotSupported
	} else if n, _ := peaker.PeakAtLeast(b, len(gzipMagic)); n == len(gzipMagic) && string(b[:len(gzipMagic)]) == gzipMagic {
		// "gzip" format. https://tools.ietf.org/html/rfc1952
		r.r, err = gzip.NewReader(r.r)
	} else if n, _ := peaker.PeakAtLeast(b, 4); n == 4 && string(b[:4]) == string("LZIP") {
		return NotSupported
	}

	// It is not a known format. Assume no compression.

	return err
}
