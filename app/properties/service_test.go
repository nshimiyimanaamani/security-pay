package properties_test

import (
	"context"
	"fmt"
	"testing"

	//"github.com/rugwirobaker/paypack-backend/app"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	wrongID    = ""
	email      = "user@example.com"
	wrongValue = "wrong-value"
)

var uuidLength = len(uuid.New().ID())

func newService(owners map[string]properties.Owner) properties.Service {
	idp := mocks.NewIdentityProvider()
	props := mocks.NewRepository(owners)
	return properties.New(idp, props)
}

func makeOwners(owner properties.Owner) map[string]properties.Owner {
	owners := make(map[string]properties.Owner)
	owners[owner.ID] = owner
	return owners
}

func TestAddProperty(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	property := properties.Property{
		Owner:      properties.Owner{ID: owner.ID},
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	invalidProperty := properties.Property{
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	emptyDue := properties.Property{
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		RecordedBy: uuid.New().ID(),
	}

	withUnsavedOwner := properties.Property{
		Owner:      properties.Owner{ID: uuid.New().ID()},
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{"add valid property", property, nil},
		{"add invalid property", invalidProperty, properties.ErrInvalidEntity},
		{"add property with empty montly due", emptyDue, properties.ErrInvalidEntity},
		{"add with unsaved owner", withUnsavedOwner, properties.ErrOwnerNotFound},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.RegisterProperty(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	property := properties.Property{
		Owner:      properties.Owner{ID: owner.ID},
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	invalidProperty := properties.Property{
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	emptyDue := properties.Property{
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		RecordedBy: uuid.New().ID(),
	}

	ctx := context.Background()
	saved, err := svc.RegisterProperty(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

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
			desc:     "update with wrong property data",
			property: invalidProperty,
			err:      properties.ErrInvalidEntity,
		},
		{
			desc:     "update non-existant property",
			property: property,
			err:      properties.ErrPropertyNotFound,
		},
		{
			desc:     "update property with empty due",
			property: emptyDue,
			err:      properties.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.UpdateProperty(ctx, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewProperty(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	property := properties.Property{
		Owner:      properties.Owner{ID: owner.ID},
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	ctx := context.Background()
	saved, err := svc.RegisterProperty(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc     string
		identity string
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
			err:      properties.ErrPropertyNotFound,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.RetrieveProperty(ctx, tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByOwner(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		property := properties.Property{
			Owner:      owner,
			Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
			Due:        float64(1000),
			RecordedBy: uuid.New().ID(),
		}

		ctx := context.Background()
		_, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := []struct {
		desc   string
		owner  string
		token  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			owner:  owner.ID,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			owner:  owner.ID,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			owner:  owner.ID,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			owner:  owner.ID,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListPropertiesByOwner(ctx, tc.owner, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesBySector(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	property := properties.Property{
		Owner:      owner,
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := []struct {
		desc   string
		sector string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			sector: property.Address.Sector,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			sector: property.Address.Sector,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			sector: property.Address.Sector,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			sector: property.Address.Sector,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListPropertiesBySector(ctx, tc.sector, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByCell(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	property := properties.Property{
		Owner:      owner,
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	}

	cases := []struct {
		desc   string
		cell   string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			cell:   property.Address.Cell,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			cell:   property.Address.Cell,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			cell:   property.Address.Cell,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			cell:   property.Address.Cell,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListPropertiesByCell(ctx, tc.cell, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestListPropertiesByVillage(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	property := properties.Property{
		Owner:      owner,
		Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := []struct {
		desc    string
		village string
		offset  uint64
		limit   uint64
		size    uint64
		err     error
	}{
		{
			desc:    "list all properties",
			village: property.Address.Village,
			offset:  0,
			limit:   n,
			size:    n,
			err:     nil,
		},
		{
			desc:    "list half of the properties",
			village: property.Address.Village,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
			err:     nil,
		},
		{
			desc: "	list empty set",
			village: property.Address.Village,
			offset:  n + 1,
			limit:   n,
			size:    0,
			err:     nil,
		},
		{
			desc:    "list with zero limit",
			village: property.Address.Village,
			offset:  1,
			limit:   0,
			size:    0,
			err:     nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.ListPropertiesByVillage(ctx, tc.village, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}
