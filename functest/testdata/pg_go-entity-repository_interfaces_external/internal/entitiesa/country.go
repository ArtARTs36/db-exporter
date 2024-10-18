//go:generate mockgen -source=country.go -package=repositoriesa -destination=../repositoriesa/mock_country.go
package entitiesa

import (
	"context"
)

type CountryRepository interface {
	Get(ctx context.Context, filter *GetCountryFilter) (*Country, error)
	List(ctx context.Context, filter *ListCountryFilter) ([]*Country, error)
	Create(ctx context.Context, country *Country) (*Country, error)
	Update(ctx context.Context, country *Country) (*Country, error)
	Delete(ctx context.Context, filter *DeleteCountryFilter) (count int64, err error)
}

type ListCountryFilter struct {
	IDs []int64
}

type GetCountryFilter struct {
	ID int64
}

type DeleteCountryFilter struct {
	IDs []int64
}

type Country struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
