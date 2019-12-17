package postgres

import (
	"context"
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/app/auth"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type authRepository struct {
	*sql.DB
}

func (repo *authRepository) Retrieve(ctx context.Context, username string) (auth.Credentials, error) {
	const op errors.Op = "store/postgres/authRepository.Retrieve"

	q := `SELECT username, role, password FROM users WHERE username=$1`

	creds := auth.Credentials{}

	if err := repo.QueryRow(q, username).Scan(&creds.Username, &creds.Role, &creds.Password); err != nil {
		if err == sql.ErrNoRows {
			return creds, errors.E(op, "user not found", errors.KindNotFound)
		}
		return creds, errors.E(op, err, errors.KindUnexpected)
	}
	return creds, nil
}
