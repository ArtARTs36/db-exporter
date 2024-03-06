package ds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/shared/ds"
)

func TestAdd(t *testing.T) {
	set := ds.NewSet()

	set.Add("1")
	set.Add("1")
	set.Add("2")
	set.Add("1")
	set.Add("3")
	set.Add("3")

	assert.Equal(t, []string{
		"1", "2", "3",
	}, set.List())
}
