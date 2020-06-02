package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveDeveloper(ctx context.Context, user users.Developer) (users.Developer, error) {
	const op errors.Op = "store/postgres.userRepository.SaveDeveloper"

	q := `
		INSERT into users (
			username, 
			password, 
			role, 
			account
		) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at;
	`

	empty := users.Developer{}

	tx, err := repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})

	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	err = tx.QueryRow(q, user.Email, user.Password, user.Role, user.Account).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				tx.Rollback()
				return empty, errors.E(op, err, "user already exists", errors.KindAlreadyExists)
			case errFK:
				tx.Rollback()
				return empty, errors.E(op, err, "invalid input data: account not found", errors.KindNotFound)
			}
		}
		tx.Rollback()
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	q = `INSERT INTO developers(email, role) VALUES ($1, $2) RETURNING email`

	_, err = tx.Exec(q, user.Email, user.Role)
	if err != nil {
		empty := users.Developer{}

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
func (repo *userRepository) RetrieveDeveloper(ctx context.Context, id string) (users.Developer, error) {
	const op errors.Op = "store/postgres/userRepository.RetrieveDeveloper"

	q := `SELECT username, account, role, created_at, updated_at FROM users WHERE username=$1`

	var user = users.Developer{}

	if err := repo.QueryRow(q, id).Scan(&user.Email, &user.Account, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		empty := users.Developer{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "user not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return user, nil
}

func (repo *userRepository) ListDevelopers(ctx context.Context, offset, limit uint64) (users.DeveloperPage, error) {
	const op errors.Op = "store/postgres/userRepository.ListDevelopers"

	q := `
		SELECT 
			username, 
			account, 
			role,  
			created_at, 
			updated_at 
		FROM 
			users 
		WHERE 
			role='dev' AND account=$1
		ORDER BY username LIMIT $2 OFFSET $3
		;
	`

	var items = []users.Developer{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, creds.Account, limit, offset)
	if err != nil {
		return users.DeveloperPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := users.Developer{}

		if err := rows.Scan(&c.Email, &c.Account, &c.Role, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return users.DeveloperPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM users WHERE role='dev' AND account=$1;`

	var total uint64
	if err := repo.QueryRow(q, creds.Account).Scan(&total); err != nil {
		return users.DeveloperPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := users.DeveloperPage{
		Developers: items,
		PageMetadata: users.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *userRepository) UpdateDeveloperCreds(ctx context.Context, user users.Developer) error {
	const op errors.Op = "store/postgres.userRepository.UpdateDeveloperCreds"

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

func (repo *userRepository) DeleteDeveloper(ctx context.Context, id string) error {
	const op errors.Op = "store/postgres/userRepository.DeleteDeveloper"

	q := `DELETE FROM users WHERE username=$1`

	res, err := repo.Exec(q, id)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	if cnt == 0 {
		return errors.E(op, "user not found", errors.KindNotFound)
	}

	return nil
}
