package postgres_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rugwirobaker/paypack-backend/app/uuid"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestSaveOwner(t *testing.T) {
	repo := postgres.NewOwnerStore(db)

	defer CleanDB(t, "owners")

	new := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	cases := []struct {
		desc  string
		owner properties.Owner
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
			err:   properties.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := repo.Save(tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateOwner(t *testing.T) {
	repo := postgres.NewOwnerStore(db)

	defer CleanDB(t, "owners")

	new := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	id, err := repo.Save(new)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	new.ID = id

	cases := []struct {
		desc  string
		owner properties.Owner
		err   error
	}{
		{
			desc:  "update existing owner",
			owner: new,
			err:   nil,
		},
		{
			desc:  "update non-existant owner",
			owner: properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"},
			err:   properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := repo.Update(tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveOwner(t *testing.T) {
	repo := postgres.NewOwnerStore(db)

	defer CleanDB(t, "owners")

	new := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	id, err := repo.Save(new)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	new.ID = id

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing owner", new.ID, nil},
		{"retrieve non-existing owner", uuid.New().ID(), properties.ErrNotFound},
		{"retrieve owner with malformed id", wrongValue, properties.ErrNotFound},
	}

	for _, tc := range cases {
		_, err := repo.Retrieve(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestFindOwner(t *testing.T) {
	repo := postgres.NewOwnerStore(db)

	defer CleanDB(t, "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	id, err := repo.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	saved, err := repo.Retrieve(id)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc  string
		fname string
		lname string
		phone string
		err   error
	}{
		{
			desc:  "find existing owner",
			fname: saved.Fname,
			lname: saved.Lname,
			phone: saved.Phone,
			err:   nil,
		},
		{
			desc:  "find owner with wrong first name",
			fname: "wrong",
			lname: saved.Lname,
			phone: saved.Phone,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with wrong last name",
			fname: saved.Fname,
			lname: "wrong",
			phone: saved.Phone,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with wrong phone number",
			fname: saved.Fname,
			lname: saved.Lname,
			phone: "wrong",
			err:   properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := repo.FindOwner(tc.fname, tc.lname, tc.phone)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveAllOwners(t *testing.T) {
	repo := postgres.NewOwnerStore(db)

	defer CleanDB(t, "owners")

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Owner{
			ID:    uuid.New().ID(),
			Fname: "James ",
			Lname: "Rodriguez",
			Phone: "0784677882",
		}
		repo.Save(p)
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
		page, err := repo.RetrieveAll(tc.offset, tc.limit)
		size := uint64(len(page.Owners))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}
