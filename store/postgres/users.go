package postgres

import (
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/app/users"
)

var _ (users.Repository) = (*userRepository)(nil)

type userRepository struct {
	*sql.DB
}

// NewUserRepository ...
func NewUserRepository(db *sql.DB) users.Repository {
	return &userRepository{db}
}
