package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
)

type messageStore struct {
	db *sql.DB
}

// NewMessageStore ...
func NewMessageStore(db *sql.DB) feedback.Repository {
	return &messageStore{db}
}

func (str *messageStore) Save(ctx context.Context, msg feedback.Message) error {
	q := `INSERT INTO messages(
			id, title, body, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := str.db.Exec(q, &msg.ID, &msg.Title, &msg.Body, &msg.CreatedBy, &msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return feedback.ErrConflict
			case errInvalid, errTruncation:
				return feedback.ErrInvalidEntity
			}
		}
		return err
	}
	return nil
}

func (str *messageStore) Retrieve(ctx context.Context, id string) (feedback.Message, error) {
	q := `SELECT id, title, body, created_by, created_at, updated_at FROM messages WHERE id=$1`

	var msg = feedback.Message{}
	if err := str.db.QueryRow(q, id).Scan(&msg.ID, &msg.Title, &msg.Body, &msg.CreatedBy, &msg.CreatedAt); err != nil {
		empty := feedback.Message{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, feedback.ErrNotFound
		}
		return empty, err
	}

	return msg, nil
}
func (str *messageStore) Update(ctx context.Context, msg feedback.Message) error {
	q := `UPDATE messages SET title=$1, body=$2, updated_at=$3 WHERE id=$4`
	res, err := str.db.Exec(q, msg.Title, msg.Body, msg.UpdatedAt, msg.ID)
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

func (str *messageStore) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM messages WHERE id=$1`

	res, err := str.db.Exec(q, id)
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
