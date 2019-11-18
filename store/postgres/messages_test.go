package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/feedback"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMessage(t *testing.T) {
	repo := postgres.NewMessageStore(db)

	defer CleanDB(t, "messages")

	message := feedback.Message{
		ID: uuid.New().ID(), Title: "title",
		Body:      "body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cases := []struct {
		desc string
		msg  feedback.Message
		err  error
	}{
		{
			desc: "save new message",
			msg:  message,
			err:  nil,
		},
		{
			desc: "save message with conflicting id",
			msg:  message,
			err:  feedback.ErrConflict,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Save(ctx, tc.msg)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestUpdateMessage(t *testing.T) {
	repo := postgres.NewMessageStore(db)

	defer CleanDB(t, "messages")

	message := feedback.Message{
		ID: uuid.New().ID(), Title: "title",
		Body:      "body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()

	err := repo.Save(ctx, message)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc string
		msg  feedback.Message
		err  error
	}{
		{
			desc: "update existant message",
			msg:  message,
			err:  nil,
		},
		{
			desc: "update non-existant message",
			msg:  feedback.Message{ID: uuid.New().ID(), Title: "title", Body: "body", UpdatedAt: time.Now()},
			err:  feedback.ErrNotFound,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Update(ctx, tc.msg)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveMessage(t *testing.T) {
	repo := postgres.NewMessageStore(db)

	defer CleanDB(t, "messages")

	message := feedback.Message{
		ID: uuid.New().ID(), Title: "title",
		Body:      "body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()

	err := repo.Save(ctx, message)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing message", message.ID, nil},
		{"retrieve non-existing owner", uuid.New().ID(), feedback.ErrNotFound},
		{"retrieve owner with malformed id", wrongValue, feedback.ErrNotFound},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Retrieve(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestDeleteMessage(t *testing.T) {
	repo := postgres.NewMessageStore(db)

	defer CleanDB(t, "messages")

	message := feedback.Message{
		ID: uuid.New().ID(), Title: "title",
		Body:      "body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()

	err := repo.Save(ctx, message)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Delete(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}
