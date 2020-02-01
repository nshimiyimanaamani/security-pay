package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/feedback"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMessage(t *testing.T) {
	repo := postgres.NewMessageRepo(db)

	defer CleanDB(t)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	message := feedback.Message{
		ID:      uuid.New().ID(),
		Title:   "title",
		Body:    "body",
		Creator: owner.Phone,
	}

	const op errors.Op = "store/postgres/messageRepo.Save"

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
			err:  errors.E(op, "conflict: message already exists"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Save(ctx, &tc.msg)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}

}

func TestUpdateMessage(t *testing.T) {
	repo := postgres.NewMessageRepo(db)

	defer CleanDB(t)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	message := feedback.Message{
		ID: uuid.New().ID(), Title: "title",
		Body:    "body",
		Creator: owner.Phone,
	}

	ctx := context.Background()

	saved, err := repo.Save(ctx, &message)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/messageRepo.Update"

	cases := []struct {
		desc string
		msg  feedback.Message
		err  error
	}{
		{
			desc: "update existant message",
			msg:  *saved,
			err:  nil,
		},
		{
			desc: "update non-existant message",
			msg:  feedback.Message{ID: uuid.New().ID(), Title: "title", Body: "body"},
			err:  errors.E(op, "message not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Update(ctx, tc.msg)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveMessage(t *testing.T) {
	repo := postgres.NewMessageRepo(db)

	defer CleanDB(t)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	message := feedback.Message{
		ID:      uuid.New().ID(),
		Title:   "title",
		Body:    "body",
		Creator: owner.Phone,
	}

	ctx := context.Background()

	saved, err := repo.Save(ctx, &message)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/messageRepo.Retrieve"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existing message",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existing owner",
			id:   uuid.New().ID(),
			err:  errors.E(op, "message not found", errors.KindNotFound),
		},
		{
			desc: "retrieve owner with malformed id",
			id:   wrongValue,
			err:  errors.E(op, "message not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Retrieve(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestDeleteMessage(t *testing.T) {
	repo := postgres.NewMessageRepo(db)

	defer CleanDB(t)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	message := feedback.Message{
		ID: uuid.New().ID(), Title: "title",
		Body:    "body",
		Creator: owner.Phone,
	}

	ctx := context.Background()

	saved, err := repo.Save(ctx, &message)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "store/postgres/messageRepo.Delete"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "delete existing message",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "delete non-existing owner",
			id:   uuid.New().ID(),
			err:  errors.E(op, "message not found", errors.KindNotFound),
		},
		{
			desc: "delete message with malformed id",
			id:   wrongValue,
			err:  errors.E(op, "invalid id: malformed uuid", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Delete(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveAllMessages(t *testing.T) {
	repo := postgres.NewMessageRepo(db)

	defer CleanDB(t)

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	owner, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		message := feedback.Message{
			ID:      uuid.New().ID(),
			Title:   "title",
			Body:    "body",
			Creator: owner.Phone,
		}

		ctx := context.Background()

		_, err := repo.Save(ctx, &message)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	}

	const op errors.Op = "store/postgres/messageRepo.RetrieveAll"

	cases := map[string]struct {
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all messages": {
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of messages": {
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.RetrieveAll(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Messages))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}
