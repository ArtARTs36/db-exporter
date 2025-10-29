package laravel

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMigrationFilename(t *testing.T) {
	name := CreateMigrationFilename("create_users_table", 1)
	assert.True(t, strings.HasSuffix(name, "create_users_table.php"))
}
