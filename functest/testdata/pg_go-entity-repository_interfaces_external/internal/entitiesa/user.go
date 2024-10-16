package entitiesa

import (
	"database/sql"
	"time"
)

type UserRepository interface {
	Get(ctx context.Context, filter *GetUserFilter) (*entitiesa.User, error)
	List(ctx context.Context, filter *ListUserFilter) ([]*entitiesa.User, error)
	Create(ctx context.Context, user *entitiesa.User) (*entitiesa.User, error)
	Update(ctx context.Context, user *entitiesa.User) (*entitiesa.User, error)
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
