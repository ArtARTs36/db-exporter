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
	tableUsers = "users"
)

type PGUserRepository struct {
	db *sqlx.DB
}

type ListUserFilter struct {
	IDs []int64
}

func NewPGUserRepository(db *sqlx.DB) *PGUserRepository {
	return &PGUserRepository{db: db}
}

func (repo *PGUserRepository) List(
	ctx context.Context,
	filter *ListUserFilter,
) ([]*entities.User, error) {
	var ents []*entities.User

	query := goqu.From(tableUsers).Select()

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
			return []*entities.User{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return ents, nil
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

func (repo *PGUserRepository) Update(
	ctx context.Context,
	user *entities.User,
) (*entities.User, error) {
	query, _, err := goqu.Update(tableUsers).
		Set(user).
		Returning("*").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	var updated entities.User

	err = repo.db.GetContext(ctx, &updated, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updated, nil
}
