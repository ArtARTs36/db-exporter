package env

import (
	"testing"

	"github.com/buildkite/interpolate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInjector_Inject(t *testing.T) {
	t.Run("success, inject env", func(t *testing.T) {
		injector := &Injector{
			env: interpolate.NewMapEnv(map[string]string{
				"DSN": "dbname=test",
			}),
		}

		got, err := injector.Inject("${DSN}")
		require.NoError(t, err)

		assert.Equal(t, "dbname=test", got)
	})

	t.Run("success, expression without env", func(t *testing.T) {
		injector := &Injector{
			env: interpolate.NewMapEnv(map[string]string{
				"DSN": "dbname=test",
			}),
		}

		got, err := injector.Inject("random string")
		require.NoError(t, err)

		assert.Equal(t, "random string", got)
	})

	t.Run("failed, env not found", func(t *testing.T) {
		injector := &Injector{
			env: interpolate.NewMapEnv(map[string]string{
				"DSN": "dbname=test",
			}),
		}

		_, err := injector.Inject("${VAR}")
		require.Error(t, err)
	})
}
