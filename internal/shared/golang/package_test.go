package golang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildPackage(t *testing.T) {
	cases := []struct {
		Title string

		PkgName string
		Module  string

		Expected Package
	}{
		{
			Title: "test single package",

			PkgName: "domain",
			Module:  "github.com/artarts36/cars",

			Expected: Package{
				Name:     "domain",
				FullName: "github.com/artarts36/cars/domain",
			},
		},
		{
			Title: "test nested package",

			PkgName: "internal/domain",
			Module:  "github.com/artarts36/cars",

			Expected: Package{
				Name:     "domain",
				FullName: "github.com/artarts36/cars/internal/domain",
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			got, err := BuildPackage(tCase.PkgName, tCase.Module)
			require.NoError(t, err)
			assert.Equal(t, tCase.Expected, got)
		})
	}
}
