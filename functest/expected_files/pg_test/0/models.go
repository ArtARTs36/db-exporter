package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `db:"id"`
	Name      sql.NullString `db:"name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}
