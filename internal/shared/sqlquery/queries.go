package sqlquery

import (
	"fmt"
)

func BuildCommentOnColumn(table, column, comment string) string {
	return fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", table, column, comment)
}

func BuildDropTable(table string, useIfExists bool) string {
	ife := ""
	if useIfExists {
		ife = "IF EXISTS "
	}

	return fmt.Sprintf("DROP TABLE %s%s;", ife, table)
}
