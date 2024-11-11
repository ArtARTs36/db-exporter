package sql

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
)

func buildCreateEmptyTable(table *schema.Table, useIfNotExists bool) string {
	return fmt.Sprintf("CREATE TABLE %s%s()", ifne(useIfNotExists), table.Name.Value)
}

func ife(use bool) string {
	if use {
		return "IF EXISTS"
	}
	return ""
}

func ifne(use bool) string {
	if use {
		return "IF NOT EXISTS "
	}
	return ""
}
