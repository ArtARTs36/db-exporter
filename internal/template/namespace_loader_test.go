package template

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamespaceLoader_splitName(t *testing.T) {
	t.Run("name is empty", func(t *testing.T) {
		l := &NamespaceLoader{}

		_, _, err := l.splitName("")

		assert.Equal(t, errors.New("empty file name"), err)
	})

	t.Run("name not contains namespace", func(t *testing.T) {
		l := &NamespaceLoader{}

		_, _, err := l.splitName("without-namespace")

		assert.Equal(t, errors.New("name not contains namespace"), err)
	})

	t.Run("ok", func(t *testing.T) {
		l := &NamespaceLoader{}

		namespace, name, err := l.splitName("@local/template.html")
		require.NoError(t, err)

		assert.Equal(t, "local", namespace)
		assert.Equal(t, "template.html", name)
	})
}
