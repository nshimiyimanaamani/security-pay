package owners_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/owners/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	wrongValue = "wrong_value"
	email      = "email"
)

func newService() owners.Service {
	idp := mocks.NewIdentityProvider()
	repo := mocks.NewRepository()
	opts := &owners.Options{Idp: idp, Repo: repo}
	return owners.New(opts)
}

func TestRegisterOwner(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc  string
		owner owners.Owner
		err   error
	}{
		{
			desc:  "add valid owner",
			owner: owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"},
			err:   nil,
		},
		{
			desc:  "add invalid owner",
			owner: owners.Owner{},
			err:   owners.ErrInvalidEntity,
		},
		{
			desc:  "add owner with empty fname field",
			owner: owners.Owner{Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"},
			err:   owners.ErrInvalidEntity,
		},
		{
			desc:  "add owner with empty lname field",
			owner: owners.Owner{Fname: "James", Phone: "0784677882", Namespace: "namespace"},
			err:   owners.ErrInvalidEntity,
		},
		{
			desc:  "add owner with invalid phone number",
			owner: owners.Owner{Fname: "James", Lname: "Torredo", Phone: "77878333", Namespace: "namespace"},
			err:   owners.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Register(ctx, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestUpdateOwner(t *testing.T) {
	svc := newService()

	ctx := context.Background()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}
	owner, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	saved := owner

	cases := []struct {
		desc  string
		owner owners.Owner
		token string
		err   error
	}{
		{
			desc:  "update existing owner",
			owner: saved,

			err: nil,
		},
		{
			desc:  "add invalid owner",
			owner: owners.Owner{},

			err: owners.ErrInvalidEntity,
		},
		{
			desc:  "update non-existant owner",
			owner: owners.Owner{Fname: "james", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"},
			err:   owners.ErrNotFound,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.Update(ctx, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveOwner(t *testing.T) {
	svc := newService()

	ctx := context.Background()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}
	owner, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	saved := owner

	cases := []struct {
		desc     string
		identity string
		token    string
		err      error
	}{
		{
			desc:     "view existing property",
			identity: saved.ID,
			err:      nil,
		},
		{
			desc:     "view non-existing property",
			identity: wrongValue,

			err: owners.ErrNotFound,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Retrieve(ctx, tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestListOwners(t *testing.T) {
	svc := newService()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}
	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.Register(ctx, owner)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := []struct {
		desc   string
		token  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.List(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Owners))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestSearch(t *testing.T) {
	svc := newService()

	fname := "james"
	lname := "torredo"
	phone := "0784677882"

	owner := owners.Owner{Fname: fname, Lname: lname, Phone: phone, Namespace: "namespace"}

	ctx := context.Background()
	saved, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc  string
		owner owners.Owner
		err   error
	}{
		{
			desc:  "find existing owner",
			owner: owners.Owner{Fname: saved.Fname, Lname: saved.Lname, Phone: saved.Phone},
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
		_, err := svc.Search(ctx, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestRetrieveByPhone(t *testing.T) {
	svc := newService()

	fname := "james"
	lname := "torredo"
	phone := "0784677882"

	owner := owners.Owner{Fname: fname, Lname: lname, Phone: phone, Namespace: "namespace"}

	ctx := context.Background()
	_, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc  string
		fname string
		lname string
		phone string
		err   error
	}{
		{
			desc:  "find existing owner with valid phone number",
			phone: phone,
			err:   nil,
		},
		{
			desc:  "find owner with wrong phone number",
			phone: "wrong",
			err:   owners.ErrNotFound,
		},
	}
	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.RetrieveByPhone(ctx, tc.phone)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}
