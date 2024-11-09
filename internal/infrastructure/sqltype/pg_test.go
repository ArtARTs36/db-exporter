package sqltype

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapPGType(t *testing.T) {
	assert.Equal(t, PGText, MapPGType("text"))
	assert.Equal(t, PGTimestampWithoutTZ, MapPGType("timestamp without time zone"))
}
