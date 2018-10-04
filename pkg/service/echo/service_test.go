package echo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	result := Echo("world")
	assert.Equal(t, "Echo world", result)
}
