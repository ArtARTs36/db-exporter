package repositoriesa

import (
	"github.com/jmoiron/sqlx"

	"github.com/project/internal/entitiesa"
)

type Container struct {
	CountryRepository entitiesa.CountryRepository
	UserRepository    entitiesa.UserRepository
}

func NewContainer(db *sqlx.DB) *Container {
	return &Container{
		CountryRepository: NewPGCountryRepository(db),
		UserRepository:    NewPGUserRepository(db),
	}
}
