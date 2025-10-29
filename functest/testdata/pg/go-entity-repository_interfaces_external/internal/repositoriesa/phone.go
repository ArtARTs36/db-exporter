//go:generate mockgen -source=phone.go -package=repositoriesa -destination=mock_phone.go
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
	tablePhones = "phones"
)

type PGPhoneRepository struct {
	db *sqlx.DB
}

func NewPGPhoneRepository(db *sqlx.DB) *PGPhoneRepository {
	return &PGPhoneRepository{db: db}
}

func (repo *PGPhoneRepository) Get(
	ctx context.Context,
	filter *entitiesa.GetPhoneFilter,
) (*entitiesa.Phone, error) {
	var ent entitiesa.Phone

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
	filter *entitiesa.ListPhoneFilter,
) ([]*entitiesa.Phone, error) {
	var ents []*entitiesa.Phone

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
			return []*entitiesa.Phone{}, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return ents, nil
}

func (repo *PGPhoneRepository) Create(
	ctx context.Context,
	phone *entitiesa.Phone,
) (*entitiesa.Phone, error) {
	query, _, err := goqu.Insert(tablePhones).Rows(phone).Returning("*").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var created entitiesa.Phone

	err = repo.db.GetContext(ctx, &created, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &created, nil
}

func (repo *PGPhoneRepository) Update(
	ctx context.Context,
	phone *entitiesa.Phone,
) (*entitiesa.Phone, error) {
	query, _, err := goqu.Update(tablePhones).
		Set(phone).
		Returning("*").
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	var updated entitiesa.Phone

	err = repo.db.GetContext(ctx, &updated, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updated, nil
}

func (repo *PGPhoneRepository) Delete(
	ctx context.Context,
	filter *entitiesa.DeletePhoneFilter,
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
