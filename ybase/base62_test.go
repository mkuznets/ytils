package ybase

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase62(t *testing.T) {
	b := []byte("hello world")
	assert.Equal(t, "aaWF93RVY4AwqvW", EncodeBase62(b))
	assert.Equal(t, b, DecodeBase62(EncodeBase62(b)))
}
