package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

func (repo *userRepository) SaveAgent(ctx context.Context, user users.Agent) (users.Agent, error) {
	const op errors.Op = "store/postgres/userRepository.SaveAgent"

	q := `
		INSERT into users (
			username, 
			password, 
			role, 
			account, 
			created_at, 
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING username;
	`

	empty := users.Agent{}

	tx, err := repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})

	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	_, err = tx.Exec(q, user.Telephone, user.Password, user.Role, user.Account, user.CreatedAt, user.UpdatedAt)

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

	q = `
		INSERT INTO agents (
			telephone, 
			role, 
			first_name, last_name,
			cell, sector, village
		) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING telephone;`

	_, err = tx.Exec(q, user.Telephone, user.Role, user.FirstName, user.LastName, user.Cell, user.Sector, user.Village)

	if err != nil {
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

func (repo *userRepository) RetrieveAgent(ctx context.Context, id string) (users.Agent, error) {
	const op errors.Op = "store/postgres/userRepository.RetrieveAgent"

	var user = users.Agent{}

	q := `SELECT username, account, role, created_at, updated_at FROM users WHERE username=$1`

	if err := repo.QueryRow(q, id).Scan(&user.Telephone, &user.Account, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		empty := users.Agent{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "user not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	q = `SELECT first_name, last_name, cell, sector, village FROM agents WHERE telephone=$1`

	if err := repo.QueryRow(q, user.Telephone).Scan(&user.FirstName, &user.LastName, &user.Cell, &user.Sector, &user.Village); err != nil {
		empty := users.Agent{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "user not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return user, nil
}

func (repo *userRepository) ListAgents(ctx context.Context, offset, limit uint64) (users.AgentPage, error) {
	const op errors.Op = "store/postgres/userRepository.ListDevelopers"

	q := `
		SELECT 
			users.username, 
			users.account, 
			users.role,  
			users.created_at, 
			users.updated_at,
			agents.cell,
			agents.sector,
			agents.village
		FROM 
			users 
		INNER JOIN 
			agents ON users.username=agents.telephone
		WHERE 
			users.role='min' 
		ORDER BY users.username LIMIT $1 OFFSET $2;
	`

	var items = []users.Agent{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return users.AgentPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := users.Agent{}

		if err := rows.Scan(&c.Telephone, &c.Account, &c.Role, &c.CreatedAt, &c.UpdatedAt, &c.Cell, &c.Sector, &c.Village); err != nil {
			return users.AgentPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM users WHERE role='min';`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return users.AgentPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := users.AgentPage{
		Agents: items,
		PageMetadata: users.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *userRepository) UpdateAgentDetails(ctx context.Context, user users.Agent) error {
	const op errors.Op = "store/postgres/userRepository.UpdateAgentDetails"

	q := `
		UPDATE agents SET 
			first_name=$1, 
			last_name=$2,
			cell=$3,
			sector=$4,
			village=$5
		WHERE telephone=$6`

	res, err := repo.Exec(q, user.FirstName, user.LastName, user.Cell, user.Sector, user.Village, user.Telephone)

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

func (repo *userRepository) UpdateAgentCreds(ctx context.Context, user users.Agent) error {
	const op errors.Op = "store/postgres/userRepository.UpdateAgentCreds"

	q := `UPDATE users SET password=$1, updated_at=$2 WHERE username=$3`

	res, err := repo.Exec(q, user.Password, user.UpdatedAt, user.Telephone)

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
func (repo *userRepository) DeleteAgent(ctx context.Context, id string) error {
	const op errors.Op = "store/postgres/userRepository.DeleteAgent"

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
