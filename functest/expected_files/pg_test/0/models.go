package models

import (
	"database/sql"
)

type User struct {
	ID   int64          `db:"id"`
	Name sql.NullString `db:"name"`
}
