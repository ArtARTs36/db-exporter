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
	IDs []int64
}

type GetUserFilter struct {
	ID int64
}

type DeleteUserFilter struct {
	IDs []int64
}

type User struct {
	ID          int64           `db:"id"`
	Name        string          `db:"name"`
	CountryID   sql.NullInt64   `db:"country_id"`
	Balance     float64         `db:"balance"`
	PrevBalance sql.NullFloat64 `db:"prev_balance"`
	Phone       sql.NullString  `db:"phone"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}
