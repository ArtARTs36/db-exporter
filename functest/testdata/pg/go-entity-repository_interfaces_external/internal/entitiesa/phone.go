//go:generate mockgen -source=phone.go -package=repositoriesa -destination=../repositoriesa/mock_phone.go
package entitiesa

import (
	"context"
)

type PhoneRepository interface {
	Get(ctx context.Context, filter *GetPhoneFilter) (*Phone, error)
	List(ctx context.Context, filter *ListPhoneFilter) ([]*Phone, error)
	Create(ctx context.Context, phone *Phone) (*Phone, error)
	Update(ctx context.Context, phone *Phone) (*Phone, error)
	Delete(ctx context.Context, filter *DeletePhoneFilter) (count int64, err error)
}

type ListPhoneFilter struct {
	UserIDs []int
	Numbers []string
}

type GetPhoneFilter struct {
	UserID int
	Number string
}

type DeletePhoneFilter struct {
	UserIDs []int
	Numbers []string
}

type Phone struct {
	UserID int    `db:"user_id"`
	Number string `db:"number"`
}
