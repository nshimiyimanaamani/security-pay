package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveProperty(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	new := properties.Property{
		ID:         nanoid.New(nil).ID(),
		Owner:      properties.Owner{ID: owner.ID},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}

	invalid := properties.Property{
		ID:        wrongValue,
		Owner:     properties.Owner{ID: "invalid"},
		Namespace: account.ID,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	const op errors.Op = "store/postgres/propertiesStore.Save"

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{
			desc:     "save new property",
			property: new,
			err:      nil,
		},

		{
			desc:     "save property with conflicting id",
			property: new,
			err:      errors.E(op, "property already exists", errors.KindAlreadyExists),
		},
		{
			desc:     "save property with invalid owner id",
			property: invalid,
			err:      errors.E(op, "owner not found", errors.KindNotFound),
		},
	}
	for _, tc := range cases {
		ctx := context.Background()
		_, err := props.Save(ctx, tc.property)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}
	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		ForRent:    true,
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}

	saved := saveProperty(t, db, property)

	const op errors.Op = "store/postgres/propertiesStore.Update"

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{
			desc:     "update existing property",
			property: saved,
			err:      nil,
		},
		{
			desc: "update non existant property",
			property: properties.Property{
				ID:        nanoid.New(nil).ID(),
				Owner:     properties.Owner{ID: uuid.New().ID()},
				Namespace: account.ID,
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: errors.E(op, "property not found", errors.KindNotFound),
		},
		{
			desc: "update property with invalid owner",
			property: properties.Property{
				ID:        nanoid.New(nil).ID(),
				Owner:     properties.Owner{ID: wrongValue},
				Namespace: account.ID,
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: errors.E(op, "invalid property", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := props.Update(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got '%v'\n", tc.desc, tc.err, err))
	}

}

func TestDelete(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		ForRent:    true,
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}
	property = saveProperty(t, db, property)

	const op errors.Op = "store/postgres/propertiesStore.Delete"

	cases := []struct {
		desc string
		uid  string
		err  error
	}{
		{
			desc: "update existing property",
			uid:  property.ID,
			err:  nil,
		},
		{
			desc: "update non existant property",
			uid:  "invalid",
			err:  errors.E(op, "property not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := props.Delete(ctx, tc.uid)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got '%v'\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveByID(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: owner.ID},
		Address: properties.Address{
			Sector:  "Gasabo",
			Cell:    "Kanserege",
			Village: "RukiriII",
		},
		Namespace:  account.ID,
		Due:        float64(1000),
		RecordedBy: agent.Telephone,
		Occupied:   true,
	}

	ctx := context.Background()
	sp, _ := props.Save(ctx, property)

	const op errors.Op = "store/postgres/propertiesStore.RetrieveByID"

	cases := []struct {
		desc  string
		id    string
		names string
		err   error
	}{
		{
			desc: "retrieve existing property",
			id:   sp.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existing property",
			id:   nanoid.New(nil).ID(),
			err:  errors.E(op, "property not found", errors.KindNotFound),
		},
		{
			desc: "retrieve with malformed id",
			id:   wrongValue,
			err:  errors.E(op, "property not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := props.RetrieveByID(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %v got '%v'\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveByOwner(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	account = saveAccount(t, db, account)

	// creds := auth.Credentials{
	// 	Username: "username",
	// 	Password: "password",
	// 	Role:     auth.Dev,
	// 	Account:  account.ID,
	// }

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	sector := "Nyarugenge"
	cell := "Kacyiru"
	village := "Kanserege"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := map[string]struct {
		owner  string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all properties with existing owner": {
			owner:  owner.ID,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of properties with existing owner": {
			owner:  owner.ID,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve properties with non-existing owner": {
			owner:  uuid.New().ID(),
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		// ctx = auth.SetECredetialsInContext(ctx, &creds)
		page, err := props.RetrieveByOwner(ctx, tc.owner, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestRetrieveBySector(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	creds := &auth.Credentials{
		Username: "username",
		Password: "password",
		Role:     auth.Dev,
		Account:  account.ID,
	}

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	sector := "Gasabo"
	cell := "Kacyiru"
	village := "Shambo"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}
	cases := map[string]struct {
		sector string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
		names  string
	}{
		"retrieve all properties with existing sector": {
			sector: sector,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of properties with existing sector": {
			sector: sector,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve properties with non-existing sector": {
			sector: wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		ctx = auth.SetECredetialsInContext(ctx, creds)
		page, err := props.RetrieveBySector(ctx, tc.sector, tc.offset, tc.limit, tc.names)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestRetrieveByCell(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	creds := &auth.Credentials{
		Username: "username",
		Password: "password",
		Role:     auth.Dev,
		Account:  account.ID,
	}

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	sector := "Gasate"
	cell := "Gasaka"
	village := "Shami"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	}
	cases := map[string]struct {
		cell   string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
		names  string
	}{
		"retrieve all properties with existing cell": {
			cell:   cell,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of properties with existing cell": {
			cell:   cell,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve properties with non-existing cell": {
			cell:   wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		ctx = auth.SetECredetialsInContext(ctx, creds)
		page, err := props.RetrieveByCell(ctx, tc.cell, tc.offset, tc.limit, tc.names)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestRetrieveByVillage(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	creds := &auth.Credentials{
		Username: "username",
		Password: "password",
		Role:     auth.Dev,
		Account:  account.ID,
	}

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	}

	cases := map[string]struct {
		village string
		offset  uint64
		limit   uint64
		size    uint64
		total   uint64
		names   string
	}{
		"retrieve all properties with existing village": {
			village: village,
			offset:  0,
			limit:   n,
			size:    n,
			total:   n,
		},
		"retrieve subset of properties with existing village": {
			village: village,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
			total:   n,
		},
		"retrieve properties with non-existing village": {
			village: wrongValue,
			offset:  0,
			limit:   n,
			size:    0,
			total:   0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		ctx = auth.SetECredetialsInContext(ctx, creds)
		page, err := props.RetrieveByVillage(ctx, tc.village, tc.offset, tc.limit, tc.names)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}

func TestRetrieveByRecorder(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	creds := &auth.Credentials{
		Username: "username",
		Password: "password",
		Role:     auth.Dev,
		Account:  account.ID,
	}

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	}

	cases := map[string]struct {
		user   string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all properties with existing village": {
			user:   agent.Telephone,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of properties with existing village": {
			user:   agent.Telephone,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve properties with non-existing village": {
			user:   wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		ctx = auth.SetECredetialsInContext(ctx, creds)
		page, err := props.RetrieveByRecorder(ctx, tc.user, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}
}
