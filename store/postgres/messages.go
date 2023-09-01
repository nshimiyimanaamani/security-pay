package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/feedback"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

type messageRepo struct {
	*sql.DB
}

// NewMessageRepo ...
func NewMessageRepo(db *sql.DB) feedback.Repository {
	return &messageRepo{db}
}

func (repo *messageRepo) Save(ctx context.Context, msg *feedback.Message) (*feedback.Message, error) {
	const op errors.Op = "store/postgres/messageRepo.Save"

	q := `INSERT INTO messages(id, title, body, creator) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at;`

	err := repo.QueryRow(q, &msg.ID, &msg.Title, &msg.Body, &msg.Creator).Scan(&msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return nil, errors.E(op, "conflict: message already exists")
			case errInvalid, errTruncation:
				return nil, errors.E(op, "invalid message: wrong uuid format", errors.KindBadRequest)
			}
		}
		return nil, err
	}
	return msg, nil
}

func (repo *messageRepo) Retrieve(ctx context.Context, id string) (feedback.Message, error) {
	const op errors.Op = "store/postgres/messageRepo.Retrieve"

	q := `
		SELECT 
			messages.id, 
			messages.title, 
			messages.body, 
			messages.creator, 
			messages.created_at, 
			messages.updated_at,
		CONCAT (
			owners.fname, ' ', owners.lname) AS name
		FROM 
			messages
		INNER JOIN 
			owners ON messages.creator=owners.phone
		WHERE messages.id=$1`

	var msg = feedback.Message{}

	err := repo.QueryRow(q, id).Scan(&msg.ID, &msg.Title, &msg.Body, &msg.Creator, &msg.CreatedAt, &msg.UpdatedAt, &msg.DisplayName)

	if err != nil {
		empty := feedback.Message{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "message not found", errors.KindNotFound)
		}
		return empty, err
	}

	return msg, nil
}
func (repo *messageRepo) Update(ctx context.Context, msg feedback.Message) error {
	const op errors.Op = "store/postgres/messageRepo.Update"

	q := `UPDATE messages SET title=$1, body=$2 WHERE id=$3`
	res, err := repo.Exec(q, msg.Title, msg.Body, msg.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.E(op, "invalid message", errors.KindBadRequest)
			}
		}
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.E(op, "message not found", errors.KindNotFound)
	}
	return nil
}

func (repo *messageRepo) Delete(ctx context.Context, id string) error {
	const op errors.Op = "store/postgres/messageRepo.Delete"

	q := `DELETE FROM messages WHERE id=$1`

	res, err := repo.Exec(q, id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.E(op, "invalid id: malformed uuid", errors.KindBadRequest)
			}
		}
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.E(op, "message not found", errors.KindNotFound)
	}
	return nil
}

func (repo *messageRepo) RetrieveAll(ctx context.Context, offset, limit uint64) (feedback.MessagePage, error) {
	const op errors.Op = "store/postgres/messageRepo.RetrieveAll"

	q := `
		SELECT 
			messages.id,
			messages.title,
			messages.body,
			messages.creator,
			messages.created_at,
			messages.updated_at,
		CONCAT (
			owners.fname, ' ', owners.lname) AS name
		FROM 
			messages
		INNER JOIN 
			owners ON messages.creator=owners.phone
		ORDER BY 
			created_at LIMIT $1 OFFSET $2;
	`

	var items = []feedback.Message{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return feedback.MessagePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		c := feedback.Message{}

		err := rows.Scan(&c.ID, &c.Title, &c.Body, &c.Creator, &c.CreatedAt, &c.UpdatedAt, &c.DisplayName)

		if err != nil {
			return feedback.MessagePage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM messages;`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return feedback.MessagePage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := feedback.MessagePage{
		Messages: items,
		PageMetadata: feedback.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}
