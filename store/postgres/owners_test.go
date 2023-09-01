package postgres_test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/core/owners"
	"github.com/nshimiyimanaamani/paypack-backend/core/uuid"

	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterBytes   = "0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func TestSaveOwner(t *testing.T) {
	repo := postgres.NewOwnerRepo(db)

	defer CleanDB(t, db)

	new := owners.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}

	cases := []struct {
		desc  string
		owner owners.Owner
		err   error
	}{
		{
			desc:  "save new owner",
			owner: new,
			err:   nil,
		},
		{
			desc:  "save owner with conflicting id",
			owner: new,
			err:   owners.ErrConflict,
		},
		{
			desc:  "save owner with invalid data",
			owner: owners.Owner{ID: "invalid", Fname: "rugwiro", Lname: "james", Phone: "0784677882"},
			err:   owners.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Save(ctx, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestUpdateOwner(t *testing.T) {
	repo := postgres.NewOwnerRepo(db)

	defer CleanDB(t, db)

	owner := owners.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}

	ctx := context.Background()
	saved, err := repo.Save(ctx, owner)

	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc  string
		owner owners.Owner
		err   error
	}{
		{
			desc:  "update existing owner",
			owner: saved,
			err:   nil,
		},
		{
			desc:  "update non-existant owner",
			owner: owners.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"},
			err:   owners.ErrNotFound,
		},
		{
			desc:  "update owner with invalid data",
			owner: owners.Owner{ID: "invalid", Fname: "rugwiro", Lname: "james", Phone: "0784677882"},
			err:   owners.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := repo.Update(ctx, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveOwner(t *testing.T) {
	repo := postgres.NewOwnerRepo(db)

	defer CleanDB(t, db)

	owner := owners.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	ctx := context.Background()

	saved, err := repo.Save(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing owner", saved.ID, nil},
		{"retrieve non-existing owner", uuid.New().ID(), owners.ErrNotFound},
		{"retrieve owner with malformed id", wrongValue, owners.ErrNotFound},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Retrieve(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestSearch(t *testing.T) {
	repo := postgres.NewOwnerRepo(db)

	defer CleanDB(t, db)

	owner := owners.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}

	ctx := context.Background()

	saved, err := repo.Save(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc  string
		owner owners.Owner
		err   error
	}{
		{
			desc:  "find existing owner",
			owner: saved,
			err:   nil,
		},
		{
			desc:  "find owner with wrong first name",
			owner: owners.Owner{Fname: "wrong_name", Lname: saved.Lname, Phone: saved.Phone},
			err:   owners.ErrNotFound,
		},
		{
			desc:  "find owner with wrong last name",
			owner: owners.Owner{Fname: saved.Fname, Lname: "wrong_name", Phone: saved.Phone},
			err:   owners.ErrNotFound,
		},
		{
			desc:  "find owner with wrong phone number",
			owner: owners.Owner{Fname: saved.Fname, Lname: saved.Lname, Phone: "wrong_number"},
			err:   owners.ErrNotFound,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Search(ctx, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveAllOwners(t *testing.T) {
	repo := postgres.NewOwnerRepo(db)

	defer CleanDB(t, db)

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := owners.Owner{
			ID:    uuid.New().ID(),
			Fname: "James ",
			Lname: "Rodriguez",
			Phone: random(15),
		}
		ctx := context.Background()
		_, err := repo.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := map[string]struct {
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all owners": {
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of all owners": {
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()

		page, err := repo.RetrieveAll(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Owners))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestRetrieveOwnerByPhone(t *testing.T) {
	repo := postgres.NewOwnerRepo(db)

	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "paypack.developers",
		Name:          "developers",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}
	account = saveAccount(t, db, account)

	new := owners.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}

	ctx := context.Background()
	_, err := repo.Save(ctx, new)

	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc  string
		phone string
		err   error
	}{
		{"retrieve existing owner", new.Phone, nil},
		{"retrieve non-existing owner", "0785460022", owners.ErrNotFound},
		{"retrieve owner with malformed id", wrongValue, owners.ErrNotFound},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.RetrieveByPhone(ctx, tc.phone)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func random(n int) string {

	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
