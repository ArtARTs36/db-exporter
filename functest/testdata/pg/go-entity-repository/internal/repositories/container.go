package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Container struct {
	PGUserRepository  *PGUserRepository
	PGPhoneRepository *PGPhoneRepository
}

func NewContainer(db *sqlx.DB) *Container {
	return &Container{
		PGUserRepository:  NewPGUserRepository(db),
		PGPhoneRepository: NewPGPhoneRepository(db),
	}
}
