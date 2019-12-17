package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveManager(ctx context.Context, user users.Manager) (users.Manager, error) {
	const op errors.Op = "store/postgres/usersRepository.SaveManager"

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
		empty := users.Manager{}
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

	q = `INSERT INTO admins(email, role, cell) VALUES ($1, $2) RETURNING email`

	_, err = repo.Exec(q, user.Email, user.Role, user.Cell)
	if err != nil {
		empty := users.Manager{}
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

func (repo *userRepository) RetrieveManager(ctx context.Context, id string) (users.Manager, error) {
	const op errors.Op = "store/postgres/usersRepository.RetrieveManager"

	q := `
		SELECT 
			users.username, 
			users.account, 
			users.role, 
			users.created_at, 
			users.update_at 
			managers.cell
		FROM 
			users WHERE username=$1
		INNER JOIN
			managers ON users.username=managers.email ;
	`

	var user = users.Manager{}

	if err := repo.QueryRow(q, id).Scan(&user.Email, &user.Account, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.Cell); err != nil {
		empty := users.Manager{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "user not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return user, nil
}

func (repo *userRepository) ListManagers(ctx context.Context, offset, limit uint64) (users.ManagerPage, error) {
	const op errors.Op = "store/postgres/usersRepository.ListManagers"

	q := `
		SELECT 
			users.username, 
			users.account, 
			users.role,  
			users.created_at, 
			users.updated_at,
			managers.cell
		FROM users ORDER BY 
			username LIMIT $1 OFFSET $2
		INNER JOIN 
			managers ON users.username=managers.email
		WHERE users.role=3;
	`

	var items = []users.Manager{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return users.ManagerPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := users.Manager{}

		if err := rows.Scan(&c.Email, &c.Account, &c.Role, &c.CreatedAt, &c.UpdatedAt, &c.Cell); err != nil {
			return users.ManagerPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM users WHERE role=3;`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return users.ManagerPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := users.ManagerPage{
		Managers: items,
		PageMetadata: users.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *userRepository) UpdateManagerCreds(ctx context.Context, user users.Manager) error {
	const op errors.Op = "store/postgres/usersRepository.UpdateManager"

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
