package sqlquery

import (
	"fmt"
)

func BuildCommentOnColumn(table, column, comment string) string {
	return fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", table, column, comment)
}

func BuildDropTable(table string) string {
	return fmt.Sprintf("DROP TABLE %s;", table)
}
