package postgres_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveProperty(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: nanoid.New(nil).ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	id, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	owner.ID = id

	new := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	invalid := properties.Property{
		ID:    wrongValue,
		Owner: "invalid",
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
			err:      properties.ErrNotFound,
		},
	}
	for _, tc := range cases {
		_, err := props.Save(tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateProperty(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: nanoid.New(nil).ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	id, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	owner.ID = id

	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	prid, _ := props.Save(property)

	property.ID = prid

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{
			desc:     "update existing property",
			property: property,
			err:      nil,
		},
		{
			desc: "update non existant property",
			property: properties.Property{
				ID:    nanoid.New(nil).ID(),
				Owner: uuid.New().ID(),
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: properties.ErrNotFound,
		},
		{
			desc: "udpate property with invalid owner id",
			property: properties.Property{
				ID:    nanoid.New(nil).ID(),
				Owner: wrongValue,
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: properties.ErrNotFound,
		},
		{
			desc: "udpate property with invalid owner",
			property: properties.Property{
				ID:    nanoid.New(nil).ID(),
				Owner: wrongValue,
				Address: properties.Address{
					Sector:  "Remera",
					Cell:    "Gishushu",
					Village: "Ingabo",
				},
				Due: float64(1000),
			},
			err: properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := props.UpdateProperty(tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestRetrieveByID(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	oid, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner.ID = oid
	property := properties.Property{
		ID:    nanoid.New(nil).ID(),
		Owner: owner.ID,
		Address: properties.Address{
			Sector:  "Gasabo",
			Cell:    "Kanserege",
			Village: "RukiriII",
		},
		Due: float64(1000),
	}

	pid, _ := props.Save(property)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{"retrieve existing property", pid, nil},
		{"retrieve non-existing property", nanoid.New(nil).ID(), properties.ErrNotFound},
		{"retrieve with malformed id", wrongValue, properties.ErrNotFound},
	}

	for _, tc := range cases {
		_, err := props.RetrieveByID(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}
func TestRetrieveByOwner(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	oid, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner.ID = oid
	sector := "Nyarugenge"
	cell := "Kacyiru"
	village := "Kanserege"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: owner.ID,
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due: float64(1000),
		}

		props.Save(p)
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
		page, err := props.RetrieveByOwner(tc.owner, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestRetrieveBySector(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	oid, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner.ID = oid
	sector := "Gasabo"
	cell := "Kacyiru"
	village := "Shambo"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: owner.ID,
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due: float64(1000),
		}

		props.Save(p)
	}
	cases := map[string]struct {
		sector string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all properties with existing owner": {
			sector: sector,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of properties with existing owner": {
			sector: sector,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve properties with non-existing owner": {
			sector: wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
	}

	for desc, tc := range cases {
		page, err := props.RetrieveBySector(tc.sector, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestRetrieveByCell(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	oid, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner.ID = oid
	sector := "Gasate"
	cell := "Gasaka"
	village := "Shami"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: owner.ID,
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due: float64(1000),
		}

		props.Save(p)
	}
	cases := map[string]struct {
		cell   string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"retrieve all properties with existing owner": {
			cell:   cell,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of properties with existing owner": {
			cell:   cell,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve properties with non-existing owner": {
			cell:   wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
	}

	for desc, tc := range cases {
		page, err := props.RetrieveByCell(tc.cell, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestRetrieveByVillage(t *testing.T) {
	props := postgres.NewPropertyStore(db)
	ows := postgres.NewOwnerStore(db)

	defer CleanDB(t, "properties", "owners")

	owner := properties.Owner{ID: uuid.New().ID(), Fname: "rugwiro", Lname: "james", Phone: "0784677882"}
	oid, err := ows.Save(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	owner.ID = oid
	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: owner.ID,
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Due: float64(1000),
		}

		props.Save(p)
	}

	cases := map[string]struct {
		village string
		offset  uint64
		limit   uint64
		size    uint64
		total   uint64
	}{
		"retrieve all properties with existing owner": {
			village: village,
			offset:  0,
			limit:   n,
			size:    n,
			total:   n,
		},
		"retrieve subset of properties with existing owner": {
			village: village,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
			total:   n,
		},
		"retrieve properties with non-existing owner": {
			village: wrongValue,
			offset:  0,
			limit:   n,
			size:    0,
			total:   0,
		},
	}

	for desc, tc := range cases {
		page, err := props.RetrieveByVillage(tc.village, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}
