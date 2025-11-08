//go:generate mockgen -source=user.go -package=repositoriesa -destination=../repositoriesa/mock_user.go
package entitiesa

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository interface {
	Get(ctx context.Context, filter *GetUserFilter) (*User, error)
	List(ctx context.Context, filter *ListUserFilter) ([]*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, filter *DeleteUserFilter) (count int64, err error)
}

type ListUserFilter struct {
	IDs []int
}

type GetUserFilter struct {
	ID int
}

type DeleteUserFilter struct {
	IDs []int
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
