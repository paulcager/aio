package io

import (
	"strings"
	"testing"

	"io"

	"github.com/stretchr/testify/assert"
)

func TestPeaker(t *testing.T) {
	r := strings.NewReader("0123456789")
	pr := NewPeakableReader(r)

	b := make([]byte, 4)
	n, err := pr.Peak(b)
	assert.NoError(t, err)
	assert.EqualValues(t, len(b), n)
	assert.Equal(t, "0123", string(b))

	// Automatically read more if asked to peak at more.
	b = make([]byte, 8)
	n, err = pr.Peak(b)
	assert.NoError(t, err)
	assert.EqualValues(t, len(b), n)
	assert.Equal(t, "01234567", string(b))

	b = make([]byte, 12)
	n, err = pr.Peak(b)
	assert.NoError(t, err)
	assert.EqualValues(t, 10, n)
	assert.Equal(t, "0123456789", string(b[:n]))

	// Now should return EOF from underlying (but with all peaked data read).
	b = make([]byte, 12)
	n, err = pr.Peak(b)
	assert.EqualValues(t, err, io.EOF)
	assert.EqualValues(t, 10, n)
	assert.Equal(t, "0123456789", string(b[:n]))

	b = make([]byte, 7)
	n, err = pr.Read(b)
	assert.NoError(t, err)
	assert.EqualValues(t, len(b), n)
	assert.Equal(t, "0123456", string(b[:n]))

	b = make([]byte, 15)
	n, err = pr.Read(b)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, n)
	assert.Equal(t, "789", string(b[:n]))

}
