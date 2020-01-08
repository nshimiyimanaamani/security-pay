package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
)

type messageStore struct {
	*sql.DB
}

// NewMessageStore ...
func NewMessageStore(db *sql.DB) feedback.Repository {
	return &messageStore{db}
}

func (repo *messageStore) Save(ctx context.Context, msg *feedback.Message) (*feedback.Message, error) {
	q := `INSERT INTO messages(id, title, body, creator) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at;`

	err := repo.QueryRow(q, &msg.ID, &msg.Title, &msg.Body, &msg.Creator).Scan(&msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return nil, feedback.ErrConflict
			case errInvalid, errTruncation:
				return nil, feedback.ErrInvalidEntity
			}
		}
		return nil, err
	}
	return msg, nil
}

func (repo *messageStore) Retrieve(ctx context.Context, id string) (feedback.Message, error) {
	q := `SELECT id, title, body, creator, created_at, updated_at FROM messages WHERE id=$1`

	var msg = feedback.Message{}
	if err := repo.QueryRow(q, id).Scan(&msg.ID, &msg.Title, &msg.Body, &msg.Creator, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
		empty := feedback.Message{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, feedback.ErrNotFound
		}
		return empty, err
	}

	return msg, nil
}
func (repo *messageStore) Update(ctx context.Context, msg feedback.Message) error {
	q := `UPDATE messages SET title=$1, body=$2 WHERE id=$3`
	res, err := repo.Exec(q, msg.Title, msg.Body, msg.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return feedback.ErrInvalidEntity
			}
		}
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return feedback.ErrNotFound
	}
	return nil
}

func (repo *messageStore) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM messages WHERE id=$1`

	res, err := repo.Exec(q, id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return feedback.ErrInvalidEntity
			}
		}
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return feedback.ErrNotFound
	}
	return nil
}
