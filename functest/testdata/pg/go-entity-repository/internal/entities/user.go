package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int             `db:"id"`
	Name        string          `db:"name"`
	Balance     float64         `db:"balance"`
	PrevBalance sql.NullFloat64 `db:"prev_balance"`
	CreatedAt   time.Time       `db:"created_at"`
	CurrentMood string          `db:"current_mood"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
	DeletedAt   sql.NullTime    `db:"deleted_at"`
}
