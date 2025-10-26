package golang

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/artarts36/gds"
)

func TestNewStringEnumOfValues(t *testing.T) {
	enum := NewStringEnumOfValues(gds.NewString("MOOD"), []string{"ok", "good"})

	assert.Equal(t, "Mood", enum.Name.Value)
	assert.Equal(t, []*StringEnumValue{
		{Name: "MoodUndefined", Value: ""},
		{Name: "MoodOk", Value: "ok"},
		{Name: "MoodGood", Value: "good"},
	}, enum.Values)
}
