package repositories

import (
	"context"
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

func NewPGCountryRepository(db *sqlx.DB) *PGCountryRepository {
	return &PGCountryRepository{db: db}
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
