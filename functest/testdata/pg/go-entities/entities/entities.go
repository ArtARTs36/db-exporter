package entities

import (
	"database/sql"
	"time"
)

type Mood string

const (
    MoodUndefined Mood = ""
    MoodOk Mood = "ok"
    MoodHappy Mood = "happy"
)

func (e Mood) Valid() bool {
    switch e {
    case MoodUndefined:
        return true
    case MoodOk:
        return true
    case MoodHappy:
        return true
    default:
        return false
    }
}

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

type Phone struct {
	UserID int    `db:"user_id"`
	Number string `db:"number"`
}
