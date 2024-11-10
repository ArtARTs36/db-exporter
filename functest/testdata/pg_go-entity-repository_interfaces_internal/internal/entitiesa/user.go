package entitiesa

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int             `db:"id"`
	Name        string          `db:"name"`
	CountryID   sql.NullInt64   `db:"country_id"`
	Balance     float64         `db:"balance"`
	PrevBalance sql.NullFloat64 `db:"prev_balance"`
	Phone       sql.NullString  `db:"phone"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}
