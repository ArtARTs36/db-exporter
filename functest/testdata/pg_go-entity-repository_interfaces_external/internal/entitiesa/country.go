package entitiesa

type CountryRepository interface {
	Get(ctx context.Context, filter *GetCountryFilter) (*entitiesa.Country, error)
	List(ctx context.Context, filter *ListCountryFilter) ([]*entitiesa.Country, error)
	Create(ctx context.Context, country *entitiesa.Country) (*entitiesa.Country, error)
	Update(ctx context.Context, country *entitiesa.Country) (*entitiesa.Country, error)
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
