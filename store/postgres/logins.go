package postgres

import (
	"context"
	"database/sql"

	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

type authRepository struct {
	*sql.DB
}

// NewAuthRepository is a postgres implementation of auth.Repository
func NewAuthRepository(db *sql.DB) auth.Repository {
	return &authRepository{db}
}

func (repo *authRepository) Retrieve(ctx context.Context, username string) (auth.Credentials, error) {
	const op errors.Op = "store/postgres/authRepository.Retrieve"

	q := `SELECT username, account, role, password FROM users WHERE username=$1`

	creds := auth.Credentials{}

	if err := repo.QueryRow(q, username).Scan(&creds.Username, &creds.Account, &creds.Role, &creds.Password); err != nil {
		if err == sql.ErrNoRows {
			return creds, errors.E(op, "user not found: invalid username or password", errors.KindNotFound)
		}
		return creds, errors.E(op, err, errors.KindUnexpected)
	}
	return creds, nil
}
