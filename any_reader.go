// Don't use standard file detection software / libmagic as it requires >= 128 bytes to be read.
// https://en.wikipedia.org/wiki/List_of_file_signatures

// ZIP files / tar files - return concat of all contained files.

package aio

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/pierrec/lz4"
)

var NotSupported = errors.New("The file type is not yet supported")

const (
	compressMagic = "\x1f\x9d"
	gzipMagic     = "\x1f\x8b"
	lzipMagic     = "LZIP"
	bzip2Magic    = "BZh"
	xzMagic       = "\xfd7zXZ\x00"
	zlibMagic     = "\x78\x9c"
	lz4Magic      = "\x04\x22\x4d\x18"
)

type AnyReader struct {
	r       io.Reader
	decided bool
}

func NewReader(r io.Reader) *AnyReader {
	return &AnyReader{r: r}
}

func (r *AnyReader) Read(b []byte) (n int, err error) {
	if !r.decided {
		err = r.decide()
		if err != nil {
			return 0, err
		}
	}

	return r.r.Read(b)
}

func (r *AnyReader) decide() error {
	var err error
	if r.decided {
		return nil
	}

	peeker := bufio.NewReader(r.r)
	r.r = peeker
	r.decided = true

	if b, err := peeker.Peek(len(compressMagic)); err == nil && string(b) == compressMagic {
		// "compress" format. https://en.wikipedia.org/wiki/Lempel-Ziv-Welch
		return NotSupported
	} else if b, err := peeker.Peek(len(gzipMagic)); err == nil && string(b) == gzipMagic {
		// "gzip" format. https://tools.ietf.org/html/rfc1952
		r.r, err = gzip.NewReader(r.r)
	} else if b, err := peeker.Peek(len(bzip2Magic)); err == nil && string(b) == bzip2Magic {
		// "bz2" format.
		r.r = bzip2.NewReader(r.r)
	} else if b, err := peeker.Peek(len(zlibMagic)); err == nil && string(b) == zlibMagic {
		// "zlib" RFC 1950
		r.r, err = zlib.NewReader(r.r)
	} else if b, err := peeker.Peek(len(lzipMagic)); err == nil && string(b) == lzipMagic {
		return NotSupported
	} else if b, err := peeker.Peek(len(lz4Magic)); err == nil && string(b) == lz4Magic {
		r.r = lz4.NewReader(r.r)
	} else if b, err := peeker.Peek(len(xzMagic)); err == nil && string(b) == xzMagic {
		r.r = NewXZReader(r.r)
	} else {
		// It is not a known format. Assume no compression.
	}

	return err
}

// NewXZReader creates a reader that decompresses the `xz` format input.
// Note that for convenience this is done by piping the input through an invocation
// of the `xz` command.
func NewXZReader(r io.Reader) io.Reader {
	return NewPipeReader(r, "xzcat")
}

func NewGZIPReader(r io.Reader) io.Reader {
	return NewPipeReader(r, "zcat")
}

func NewPipeReader(r io.Reader, cmdName string, args ...string) io.ReadCloser {
	rpipe, wpipe := io.Pipe()

	cmd := exec.Command(cmdName, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = r
	cmd.Stdout = wpipe

	go func() {
		err := cmd.Run()
		wpipe.CloseWithError(err)
	}()

	return rpipe
}
