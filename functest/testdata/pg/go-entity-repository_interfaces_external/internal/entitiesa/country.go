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
	IDs []int
}

type GetCountryFilter struct {
	ID int
}

type DeleteCountryFilter struct {
	IDs []int
}

type Country struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
