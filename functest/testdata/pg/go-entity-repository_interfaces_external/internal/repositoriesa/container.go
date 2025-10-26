package repositoriesa

import (
	"github.com/jmoiron/sqlx"

	"github.com/project/internal/entitiesa"
)

type Container struct {
	UserRepository  entitiesa.UserRepository
	PhoneRepository entitiesa.PhoneRepository
}

func NewContainer(db *sqlx.DB) *Container {
	return &Container{
		UserRepository:  NewPGUserRepository(db),
		PhoneRepository: NewPGPhoneRepository(db),
	}
}
