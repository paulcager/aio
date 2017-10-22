package io

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"io"
)

func TestReadEmpty(t *testing.T) {
	r := NewAnyReader(strings.NewReader(""))
	b := make([]byte, 12)
	n, err := r.Read(b)
	assert.EqualValues(t, 0, n)
	assert.True(t, err == nil || err == io.EOF, "Err was %s", err)
}

func TestReadPlain(t *testing.T) {
	const str = "HelloWorld"
	r := NewAnyReader(strings.NewReader(str))
	b := make([]byte, len(str)+12)
	n, err := r.Read(b)
	assert.NoError(t,err)
	assert.EqualValues(t, len(str), n)
	assert.EqualValues(t, str, string(b[:n]))

	n, err = r.Read(b)
	assert.EqualValues(t, 0, n)
	assert.EqualValues(t, io.EOF, err)
}

func TestReadPlainShortReads(t *testing.T) {
	const str = "HelloWorld"
	r := NewAnyReader(strings.NewReader(str))
	b := make([]byte, 1)
	for i := range str {
		n, err := r.Read(b)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, n)
		assert.EqualValues(t, str[i], b[0])
	}
	n, err := r.Read(b)
	assert.EqualValues(t, 0, n)
	assert.EqualValues(t, io.EOF, err)
}

func TestReadEmptyGZIP(t *testing.T) {

}
