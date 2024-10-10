package repositories

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/project/internal/entities"
)

const (
	tableUsers = "users"
)

type PGUserRepository struct {
	db *sqlx.DB
}

func NewPGUserRepository(db *sqlx.DB) *PGUserRepository {
	return &PGUserRepository{db: db}
}

func (repo *PGUserRepository) Create(
	ctx context.Context,
	user *entities.User,
) (*entities.User, error) {
	query, _, err := goqu.Insert(tableUsers).Rows(user).Returning("*").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var created entities.User

	err = repo.db.GetContext(ctx, &created, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &created, nil
}
