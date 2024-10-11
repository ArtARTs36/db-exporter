package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/project/internal/entities"
)

const (
	tableCountries = "countries"
)

type PGCountryRepository struct {
	db *sqlx.DB
}

type ListCountryFilter struct {
	IDs []int64
}

func NewPGCountryRepository(db *sqlx.DB) *PGCountryRepository {
	return &PGCountryRepository{db: db}
}

func (repo *PGCountryRepository) List(
	ctx context.Context,
	filter *ListCountryFilter,
) ([]*entities.Country, error) {
	var ents []*entities.Country

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
			return []*entities.Country{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return ents, nil
}

func (repo *PGCountryRepository) Create(
	ctx context.Context,
	country *entities.Country,
) (*entities.Country, error) {
	query, _, err := goqu.Insert(tableCountries).Rows(country).Returning("*").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var created entities.Country

	err = repo.db.GetContext(ctx, &created, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &created, nil
}

func (repo *PGCountryRepository) Update(
	ctx context.Context,
	country *entities.Country,
) (*entities.Country, error) {
	query, _, err := goqu.Update(tableCountries).
		Set(country).
		Returning("*").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	var updated entities.Country

	err = repo.db.GetContext(ctx, &updated, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updated, nil
}
