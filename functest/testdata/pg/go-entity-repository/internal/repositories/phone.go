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
	tablePhones = "phones"
)

type PGPhoneRepository struct {
	db *sqlx.DB
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

func NewPGPhoneRepository(db *sqlx.DB) *PGPhoneRepository {
	return &PGPhoneRepository{db: db}
}

func (repo *PGPhoneRepository) Get(
	ctx context.Context,
	filter *GetPhoneFilter,
) (*entities.Phone, error) {
	var ent entities.Phone

	query := goqu.From(tablePhones).Select().Limit(1)

	if filter != nil {
		if filter.UserID > 0 {
			query = query.Where(goqu.C("user_id").Eq(filter.UserID))
		}
		if len(filter.Number) > 0 {
			query = query.Where(goqu.C("number").Eq(filter.Number))
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

func (repo *PGPhoneRepository) List(
	ctx context.Context,
	filter *ListPhoneFilter,
) ([]*entities.Phone, error) {
	var ents []*entities.Phone

	query := goqu.From(tablePhones).Select()

	if filter != nil {
		if len(filter.UserIDs) > 0 {
			query = query.Where(goqu.C("user_id").In(filter.UserIDs))
		}
		if len(filter.Numbers) > 0 {
			query = query.Where(goqu.C("number").In(filter.Numbers))
		}
	}

	q, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = repo.db.SelectContext(ctx, &ents, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*entities.Phone{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return ents, nil
}

func (repo *PGPhoneRepository) Create(
	ctx context.Context,
	phone *entities.Phone,
) (*entities.Phone, error) {
	query, _, err := goqu.Insert(tablePhones).Rows(phone).Returning("*").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var created entities.Phone

	err = repo.db.GetContext(ctx, &created, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &created, nil
}

func (repo *PGPhoneRepository) Update(
	ctx context.Context,
	phone *entities.Phone,
) (*entities.Phone, error) {
	query, _, err := goqu.Update(tablePhones).
		Set(phone).
		Returning("*").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	var updated entities.Phone

	err = repo.db.GetContext(ctx, &updated, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updated, nil
}

func (repo *PGPhoneRepository) Delete(
	ctx context.Context,
	filter *DeletePhoneFilter,
) (count int64, err error) {
	query := goqu.From(tablePhones).Delete()

	if filter != nil {
		if len(filter.UserIDs) > 0 {
			query = query.Where(goqu.C("user_id").In(filter.UserIDs))
		}
		if len(filter.Numbers) > 0 {
			query = query.Where(goqu.C("number").In(filter.Numbers))
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
