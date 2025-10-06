//go:generate mockgen -source=country.go -package=repositoriesa -destination=mock_country.go
package repositoriesa

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/project/internal/entitiesa"
)

const (
	tableCountries = "countries"
)

type CountryRepository interface {
	Get(ctx context.Context, filter *GetCountryFilter) (*entitiesa.Country, error)
	List(ctx context.Context, filter *ListCountryFilter) ([]*entitiesa.Country, error)
	Create(ctx context.Context, country *entitiesa.Country) (*entitiesa.Country, error)
	Update(ctx context.Context, country *entitiesa.Country) (*entitiesa.Country, error)
	Delete(ctx context.Context, filter *DeleteCountryFilter) (count int64, err error)
}

type PGCountryRepository struct {
	db *sqlx.DB
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

func NewPGCountryRepository(db *sqlx.DB) *PGCountryRepository {
	return &PGCountryRepository{db: db}
}

func (repo *PGCountryRepository) Get(
	ctx context.Context,
	filter *GetCountryFilter,
) (*entitiesa.Country, error) {
	var ent entitiesa.Country

	query := goqu.From(tableCountries).Select().Limit(1)

	if filter != nil {
		if filter.ID > 0 {
			query = query.Where(goqu.C("id").Eq(filter.ID))
		}
	}

	q, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = repo.db.GetContext(ctx, &ent, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &ent, nil
}

func (repo *PGCountryRepository) List(
	ctx context.Context,
	filter *ListCountryFilter,
) ([]*entitiesa.Country, error) {
	var ents []*entitiesa.Country

	query := goqu.From(tableCountries).Select()

	if filter != nil {
		if len(filter.IDs) > 0 {
			query = query.Where(goqu.C("id").In(filter.IDs))
		}
	}

	q, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = repo.db.SelectContext(ctx, &ents, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*entitiesa.Country{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return ents, nil
}

func (repo *PGCountryRepository) Create(
	ctx context.Context,
	country *entitiesa.Country,
) (*entitiesa.Country, error) {
	query, _, err := goqu.Insert(tableCountries).Rows(country).Returning("*").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var created entitiesa.Country

	err = repo.db.GetContext(ctx, &created, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &created, nil
}

func (repo *PGCountryRepository) Update(
	ctx context.Context,
	country *entitiesa.Country,
) (*entitiesa.Country, error) {
	query, _, err := goqu.Update(tableCountries).
		Set(country).
		Returning("*").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	var updated entitiesa.Country

	err = repo.db.GetContext(ctx, &updated, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updated, nil
}

func (repo *PGCountryRepository) Delete(
	ctx context.Context,
	filter *DeleteCountryFilter,
) (count int64, err error) {
	query := goqu.From(tableCountries).Delete()

	if filter != nil {
		if len(filter.IDs) > 0 {
			query = query.Where(goqu.C("id").In(filter.IDs))
		}
	}

	q, args, err := query.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := repo.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	count, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return
}
