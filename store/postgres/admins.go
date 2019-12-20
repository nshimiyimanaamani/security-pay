package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// SaveAdmin ...
func (repo *userRepository) SaveAdmin(ctx context.Context, user users.Administrator) (users.Administrator, error) {
	const op errors.Op = "store/postgres/userRepository.SaveAdmin"

	q := `
		INSERT into users (
			username, 
			password, 
			role, 
			account, 
			created_at, 
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING username
	`

	empty := users.Administrator{}

	tx, err := repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})

	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	_, err = tx.Exec(q, user.Email, user.Password, user.Role, user.Account, user.CreatedAt, user.UpdatedAt)
	if err != nil {

		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				tx.Rollback()
				return empty, errors.E(op, err, "user already exists", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				tx.Rollback()
				return empty, errors.E(op, err, "invalid account data", errors.KindBadRequest)
			}
		}
		tx.Rollback()
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	q = `INSERT INTO admins(email, role) VALUES ($1, $2) RETURNING email`

	_, err = tx.Exec(q, user.Email, user.Role)
	if err != nil {
		empty := users.Administrator{}
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				tx.Rollback()
				return empty, errors.E(op, err, "user already exists", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				tx.Rollback()
				return empty, errors.E(op, err, "invalid account data", errors.KindBadRequest)
			}
		}
		tx.Rollback()
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	tx.Commit()
	return user, nil
}

func (repo *userRepository) RetrieveAdmin(ctx context.Context, id string) (users.Administrator, error) {
	const op errors.Op = "store/postgres/userRepository.RetrieveAdmin"

	q := `SELECT username, account, role, created_at, updated_at FROM users WHERE username=$1`

	var user = users.Administrator{}

	if err := repo.QueryRow(q, id).Scan(&user.Email, &user.Account, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		empty := users.Administrator{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "user not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return user, nil
}

func (repo *userRepository) ListAdmins(ctx context.Context, offset, limit uint64) (users.AdministratorPage, error) {
	const op errors.Op = "store/postgres/userRepository.ListAdmins"

	q := `
		SELECT 
			username, 
			account, 
			role,  
			created_at, 
			updated_at 
		FROM users WHERE role='admin' ORDER BY username LIMIT $1 OFFSET $2;
	`

	var items = []users.Administrator{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return users.AdministratorPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := users.Administrator{}

		if err := rows.Scan(&c.Email, &c.Account, &c.Role, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return users.AdministratorPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM users WHERE role='admin';`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return users.AdministratorPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := users.AdministratorPage{
		Administrators: items,
		PageMetadata: users.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *userRepository) UpdateAdminCreds(ctx context.Context, user users.Administrator) error {
	const op errors.Op = "store/postgres/userRepository.UpdateAdminCreds"

	q := `UPDATE users SET password=$1, updated_at=$2 WHERE username=$3`

	res, err := repo.Exec(q, user.Password, user.UpdatedAt, user.Email)

	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	if cnt == 0 {
		return errors.E(op, err, "user not found", errors.KindNotFound)
	}
	return nil
}
