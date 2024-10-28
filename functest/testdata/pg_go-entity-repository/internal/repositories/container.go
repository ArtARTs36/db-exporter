package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Container struct {
	PGCountryRepository *PGCountryRepository
	PGUserRepository    *PGUserRepository
}

func NewContainer(db *sqlx.DB) *Container {
	return &Container{
		PGCountryRepository: NewPGCountryRepository(db),
		PGUserRepository:    NewPGUserRepository(db),
	}
}
