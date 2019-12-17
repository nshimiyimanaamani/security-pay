package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveDeveloper(ctx context.Context, user users.Developer) (users.Developer, error) {
	const op errors.Op = "store/postgres/usersRepository.SaveDeveloper"

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
	_, err := repo.Exec(q, user.Email, user.Password, user.Role, user.Account, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		empty := users.Developer{}

		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, errors.E(op, "user already exists", err, errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return empty, errors.E(op, err, "invalid user data", err, errors.KindBadRequest)
			}
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	q = `INSERT INTO admins(email, role) VALUES ($1, $2) RETURNING email`

	_, err = repo.Exec(q, user.Email, user.Role)
	if err != nil {
		empty := users.Developer{}

		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, errors.E(op, "user already exists", err, errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return empty, errors.E(op, err, "invalid user data", err, errors.KindBadRequest)
			}
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	return user, nil
}
func (repo *userRepository) RetrieveDeveloper(ctx context.Context, id string) (users.Developer, error) {
	const op errors.Op = "store/postgres/usersRepository.RetrieveDeveloper"

	q := `SELECT username, account, role, created_at, update_at FROM users WHERE username=$1`

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
	const op errors.Op = "store/postgres/usersRepository.ListDevelopers"

	q := `
		SELECT 
			username, 
			account, 
			role,  
			created_at, 
			updated_at 
		FROM users 
			ORDER BY username LIMIT $1 OFFSET $2
		WHERE users.role=1;`

	var items = []users.Developer{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return users.DeveloperPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := users.Developer{}

		if err := rows.Scan(); err != nil {
			return users.DeveloperPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM users WHERE role=1;`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
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
	const op errors.Op = "store/postgres/usersRepository.UpdateDeveloper"

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
		return errors.E(op, "user not found", errors.KindNotFound)
	}
	return nil
}
