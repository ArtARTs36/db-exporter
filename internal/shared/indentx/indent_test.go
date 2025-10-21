package indentx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	indent := NewIndent(2)

	assert.Equal(t, "", indent.Curr())
	assert.Equal(t, "  ", indent.Next().Curr())
	assert.Equal(t, "    ", indent.Next().Next().Curr())
}
