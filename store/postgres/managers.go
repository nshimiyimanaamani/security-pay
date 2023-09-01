package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/users"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveManager(ctx context.Context, user users.Manager) (users.Manager, error) {
	const op errors.Op = "store/postgres.userRepository.SaveManager"

	q := `
		INSERT into users (
			username, 
			password, 
			role, 
			account
		) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at;
	`

	empty := users.Manager{}

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
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	q = `INSERT INTO managers(email, role, cell) VALUES ($1, $2, $3) RETURNING email`

	_, err = tx.Exec(q, user.Email, user.Role, user.Cell)
	if err != nil {
		empty := users.Manager{}
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				tx.Rollback()
				return empty, errors.E(op, err, "user already exists", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				tx.Rollback()
				return empty, errors.E(op, err, "invalid user data", errors.KindBadRequest)
			}
		}
		tx.Rollback()
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	tx.Commit()
	return user, nil
}

func (repo *userRepository) RetrieveManager(ctx context.Context, id string) (users.Manager, error) {
	const op errors.Op = "store/postgres/userRepository.RetrieveManager"

	q := `
		SELECT 
			users.username, 
			users.account, 
			users.role, 
			users.password,
			users.created_at, 
			users.updated_at,
			managers.cell
		FROM 
			users 
		INNER JOIN
			managers ON users.username=managers.email
		WHERE users.username=$1;
	`

	var user = users.Manager{}

	if err := repo.QueryRow(q, id).Scan(&user.Email, &user.Account, &user.Role, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Cell); err != nil {
		empty := users.Manager{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "user not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return user, nil
}

func (repo *userRepository) ListManagers(ctx context.Context, offset, limit uint64) (users.ManagerPage, error) {
	const op errors.Op = "store/postgres/userRepository.ListManagers"

	q := `
		SELECT 
			users.username, 
			users.account, 
			users.role,  
			users.created_at, 
			users.updated_at,
			managers.cell
		FROM 
			users 
		INNER JOIN 
			managers ON users.username=managers.email
		WHERE 
			users.role='basic' and account=$1
		ORDER 
			BY username LIMIT $2 OFFSET $3;
	`

	var items = []users.Manager{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, creds.Account, limit, offset)
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

	q = `SELECT COUNT(*) FROM users WHERE role='basic' AND account=$1;`

	var total uint64
	if err := repo.QueryRow(q, creds.Account).Scan(&total); err != nil {
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
	const op errors.Op = "store/postgres.userRepository.UpdateManagerCreds"

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

func (repo *userRepository) DeleteManager(ctx context.Context, id string) error {
	const op errors.Op = "store/postgres/userRepository.DeleteManager"

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
