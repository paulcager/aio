package aio

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	z, err := NewZipFlattener("testdata/test1.zip")
	require.NoError(t, err)
	require.NotNil(t, z)
	b, err := ioutil.ReadAll(z)
	assert.Equal(t, "HelloWorld", string(b))
}
