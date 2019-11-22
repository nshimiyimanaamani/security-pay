package properties_test

import (
	"fmt"
	"testing"

	//"github.com/rugwirobaker/paypack-backend/app"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	wrongID    = ""
	email      = "user@example.com"
	token      = "token"
	wrongValue = "wrong-value"
)

func newService(tokens map[string]string) properties.Service {
	auth := mocks.NewAuthBackend(tokens)
	idp := mocks.NewIdentityProvider()
	props := mocks.NewRepository()
	return properties.New(idp, props, auth)
}

func TestAddProperty(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner:   properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	invalidProperty := properties.Property{
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	emptyDue := properties.Property{
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
	}

	cases := []struct {
		desc     string
		property properties.Property
		token    string
		err      error
	}{
		{"add valid property", property, token, nil},
		{"add invalid property", invalidProperty, token, properties.ErrInvalidEntity},
		{"add property with empty montly due", emptyDue, token, properties.ErrInvalidEntity},
		{"add property with wrong user token", property, wrongValue, users.ErrUnauthorizedAccess},
	}

	for _, tc := range cases {
		_, err := svc.RegisterProperty(tc.token, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner:   properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	invalidProperty := properties.Property{
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	emptyDue := properties.Property{
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
	}

	saved, _ := svc.RegisterProperty(token, property)

	cases := []struct {
		desc     string
		property properties.Property
		token    string
		err      error
	}{
		{
			desc:     "update existing property",
			property: saved,
			token:    token,
			err:      nil,
		},
		{
			desc:     "update with wrong property data",
			property: invalidProperty,
			token:    token,
			err:      properties.ErrInvalidEntity,
		},
		{
			desc:     "update non-existant property",
			property: property,
			token:    token,
			err:      properties.ErrNotFound,
		},
		{
			desc:     "update property with empty due",
			property: emptyDue,
			token:    token,
			err:      properties.ErrInvalidEntity,
		},
		{
			desc:     "update property with wrong token",
			property: saved,
			token:    wrongValue,
			err:      users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateProperty(tc.token, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewProperty(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner:   properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	saved, _ := svc.RegisterProperty(token, property)

	cases := []struct {
		desc     string
		identity string
		token    string
		err      error
	}{
		{
			desc:     "view existing property",
			identity: saved.ID,
			token:    token,
			err:      nil,
		},
		{
			desc:     "view non-existing property",
			identity: wrongValue,
			token:    token,
			err:      properties.ErrNotFound,
		},
		{
			desc:     "view non-existing property",
			identity: wrongValue,
			token:    wrongValue,
			err:      users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		_, err := svc.RetrieveProperty(tc.token, tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByOwner(t *testing.T) {
	svc := newService(map[string]string{token: email})

	owner := properties.Owner{ID: uuid.New().ID()}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		property := properties.Property{
			Owner:   owner,
			Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
			Due:     float64(1000),
		}
		_, err := svc.RegisterProperty(token, property)
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
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			owner:  owner.ID,
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			owner:  owner.ID,
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			owner:  owner.ID,
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list properties with invalid token",
			owner:  owner.ID,
			token:  wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			err:    users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesByOwner(tc.token, tc.owner, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesBySector(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner:   properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.RegisterProperty(token, property)
	}

	cases := []struct {
		desc   string
		sector string
		token  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			sector: property.Address.Sector,
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			sector: property.Address.Sector,
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			sector: property.Address.Sector,
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			sector: property.Address.Sector,
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with invalid token",
			sector: property.Address.Sector,
			token:  wrongValue,
			offset: 1,
			limit:  0,
			size:   0,
			err:    users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesBySector(tc.token, tc.sector, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByCell(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner:   properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.RegisterProperty(token, property)
	}

	cases := []struct {
		desc   string
		cell   string
		token  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			cell:   property.Address.Cell,
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			cell:   property.Address.Cell,
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			cell:   property.Address.Cell,
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			cell:   property.Address.Cell,
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with invalid token",
			cell:   property.Address.Cell,
			token:  wrongValue,
			offset: 1,
			limit:  0,
			size:   0,
			err:    users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesByCell(tc.token, tc.cell, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestListPropertiesByVillage(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner:   properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.RegisterProperty(token, property)
	}

	cases := []struct {
		desc    string
		village string
		token   string
		offset  uint64
		limit   uint64
		size    uint64
		err     error
	}{
		{
			desc:    "list all properties",
			village: property.Address.Village,
			token:   token,
			offset:  0,
			limit:   n,
			size:    n,
			err:     nil,
		},
		{
			desc:    "list half of the properties",
			village: property.Address.Village,
			token:   token,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
			err:     nil,
		},
		{
			desc: "	list empty set",
			village: property.Address.Village,
			token:   token,
			offset:  n + 1,
			limit:   n,
			size:    0,
			err:     nil,
		},
		{
			desc:    "list with zero limit",
			village: property.Address.Village,
			token:   token,
			offset:  1,
			limit:   0,
			size:    0,
			err:     nil,
		},
		{
			desc:    "list with zero limit",
			village: property.Address.Village,
			token:   wrongValue,
			offset:  1,
			limit:   0,
			size:    0,
			err:     users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesByVillage(tc.token, tc.village, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}
