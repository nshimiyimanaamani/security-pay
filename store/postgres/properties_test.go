package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveProperty(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "properties", "owners", "users")

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	savedOwner, err := saveOwner(t, db, owner)

	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	new := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: savedOwner.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due:        float64(1000),
		RecordedBy: savedUser.ID,
		Occupied:   true,
	}

	invalid := properties.Property{
		ID:    wrongValue,
		Owner: properties.Owner{ID: "invalid"},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

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
			err:      properties.ErrConflict,
		},
		{
			desc:     "save property with invalid owner id",
			property: invalid,
			err:      properties.ErrInvalidEntity,
		},
	}
	for _, tc := range cases {
		ctx := context.Background()
		_, err := props.Save(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateProperty(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	defer CleanDB(t, "properties", "owners", "users")

	savedOwner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}

	sown, err := saveOwner(t, db, savedOwner)

	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: sown.ID},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due:        float64(1000),
		RecordedBy: savedUser.ID,
		Occupied:   true,
	}

	ctx := context.Background()
	sp, err := props.Save(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{
			desc:     "update existing property",
			property: sp,
			err:      nil,
		},
		{
			desc: "update non existant property",
			property: properties.Property{
				ID:    nanoid.New(nil).ID(),
				Owner: properties.Owner{ID: uuid.New().ID()},
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: properties.ErrPropertyNotFound,
		},
		{
			desc: "udpate property with invalid owner",
			property: properties.Property{
				ID:    nanoid.New(nil).ID(),
				Owner: properties.Owner{ID: wrongValue},
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: properties.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := props.UpdateProperty(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveByID(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "properties", "owners", "users")

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	sown, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: properties.Owner{ID: sown.ID},
		Address: properties.Address{
			Sector:  "Gasabo",
			Cell:    "Kanserege",
			Village: "RukiriII",
		},
		Due:        float64(1000),
		RecordedBy: savedUser.ID,
		Occupied:   true,
	}

	ctx := context.Background()
	sp, _ := props.Save(ctx, property)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing property", sp.ID, nil},
		{"retrieve non-existing property", nanoid.New(nil).ID(), properties.ErrPropertyNotFound},
		{"retrieve with malformed id", wrongValue, properties.ErrPropertyNotFound},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := props.RetrieveByID(ctx, tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveByOwner(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "properties", "owners", "users")

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	saved, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	sector := "Nyarugenge"
	cell := "Kacyiru"
	village := "Kanserege"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: saved.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due:        float64(1000),
			RecordedBy: savedUser.ID,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
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
		page, err := props.RetrieveByOwner(ctx, tc.owner, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestRetrieveBySector(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "properties", "owners", "users")

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	saved, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	sector := "Gasabo"
	cell := "Kacyiru"
	village := "Shambo"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: saved.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due:        float64(1000),
			RecordedBy: savedUser.ID,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}
	cases := map[string]struct {
		sector string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
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
		page, err := props.RetrieveBySector(ctx, tc.sector, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestRetrieveByCell(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "properties", "owners", "users")

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	saved, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	sector := "Gasate"
	cell := "Gasaka"
	village := "Shami"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: saved.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due:        float64(1000),
			RecordedBy: savedUser.ID,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	}
	cases := map[string]struct {
		cell   string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
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
		page, err := props.RetrieveByCell(ctx, tc.cell, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestRetrieveByVillage(t *testing.T) {
	props := postgres.NewPropertyStore(db)

	defer CleanDB(t, "properties", "owners", "users")

	user := users.User{ID: uuid.New().ID(), Email: "email", Password: "password"}
	savedUser, err := saveUser(t, db, user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	saved, err := saveOwner(t, db, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: saved.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due:        float64(1000),
			RecordedBy: savedUser.ID,
			Occupied:   true,
		}

		ctx := context.Background()
		_, err := props.Save(ctx, p)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	}

	cases := map[string]struct {
		village string
		offset  uint64
		limit   uint64
		size    uint64
		total   uint64
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
		page, err := props.RetrieveByVillage(ctx, tc.village, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}
