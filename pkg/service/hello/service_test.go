package hello

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	result := Hello("world")
	assert.Equal(t, "Hello world", result)
}
