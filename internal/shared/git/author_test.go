package git_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/artarts36/db-exporter/internal/shared/git"
)

func TestNewAuthor(t *testing.T) {
	cases := []struct {
		Identity string
		Expected *git.Author
		Err      error
	}{
		{
			Identity: "simple <simple@mail.ru>",
			Expected: &git.Author{
				Name:     "simple",
				Email:    "simple@mail.ru",
				Identity: "simple <simple@mail.ru>",
			},
		},
		{
			Identity: "github-actions[bot] <github-actions[bot]@users.noreply.github.com>",
			Expected: &git.Author{
				Name:     "github-actions[bot]",
				Email:    "github-actions[bot]@users.noreply.github.com",
				Identity: "github-actions[bot] <github-actions[bot]@users.noreply.github.com>",
			},
		},
		{
			Identity: " <github-actions[bot]@users.noreply.github.com>",
			Err:      errors.New("failed to parse author identity: name is empty"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Identity, func(t *testing.T) {
			a, err := git.NewAuthor(tCase.Identity)
			if tCase.Err != nil {
				require.Error(t, tCase.Err, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tCase.Expected, a)
		})
	}
}
