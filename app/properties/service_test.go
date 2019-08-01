package properties_test

import (
	"fmt"
	"testing"

	//"github.com/rugwirobaker/paypack-backend/app"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/rugwirobaker/paypack-backend/app/users"
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
	props := mocks.NewPropertyStore()
	owners := mocks.NewOwnerStore()
	return properties.New(idp, owners, props, auth)
}

func TestAddProperty(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
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
		_, err := svc.AddProperty(tc.token, tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	invalidProperty := properties.Property{
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due:     float64(1000),
	}

	emptyDue := properties.Property{
		Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
	}

	saved, _ := svc.AddProperty(token, property)

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
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	saved, _ := svc.AddProperty(token, property)

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
		_, err := svc.ViewProperty(tc.token, tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByOwner(t *testing.T) {
	svc := newService(map[string]string{token: email})

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}
	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}

	oid, err := svc.CreateOwner(token, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		property.Owner = oid
		_, err = svc.AddProperty(token, property)
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
			owner:  property.Owner,
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			owner:  property.Owner,
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			owner:  property.Owner,
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			owner:  property.Owner,
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list properties with invalid token",
			owner:  property.Owner,
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
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(token, property)
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
			sector: property.Sector,
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			sector: property.Sector,
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			sector: property.Sector,
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			sector: property.Sector,
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with invalid token",
			sector: property.Sector,
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
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(token, property)
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
			cell:   property.Cell,
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			cell:   property.Cell,
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			cell:   property.Cell,
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			cell:   property.Cell,
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with invalid token",
			cell:   property.Cell,
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
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(token, property)
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
			village: property.Village,
			token:   token,
			offset:  0,
			limit:   n,
			size:    n,
			err:     nil,
		},
		{
			desc:    "list half of the properties",
			village: property.Village,
			token:   token,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
			err:     nil,
		},
		{
			desc: "	list empty set",
			village: property.Village,
			token:   token,
			offset:  n + 1,
			limit:   n,
			size:    0,
			err:     nil,
		},
		{
			desc:    "list with zero limit",
			village: property.Village,
			token:   token,
			offset:  1,
			limit:   0,
			size:    0,
			err:     nil,
		},
		{
			desc:    "list with zero limit",
			village: property.Village,
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

func TestCreateOwner(t *testing.T) {
	svc := newService(map[string]string{token: email})

	cases := []struct {
		desc  string
		owner properties.Owner
		token string
		err   error
	}{
		{
			desc:  "add valid owner",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"},
			token: token,
			err:   nil,
		},
		{
			desc:  "add invalid owner",
			owner: properties.Owner{},
			token: token,
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with empty fname field",
			owner: properties.Owner{Lname: "Torredo", Phone: "0784677882"},
			token: token,
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with empty lname field",
			owner: properties.Owner{Fname: "James", Phone: "0784677882"},
			token: token,
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with invalid phone number",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "77878333"},
			token: token,
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with invalid token",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"},
			token: wrongValue,
			err:   users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		_, err := svc.CreateOwner(tc.token, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateOwner(t *testing.T) {
	svc := newService(map[string]string{token: email})

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	owner.ID, _ = svc.CreateOwner(token, owner)

	saved := owner

	cases := []struct {
		desc  string
		owner properties.Owner
		token string
		err   error
	}{
		{
			desc:  "update existing owner",
			owner: saved,
			token: token,
			err:   nil,
		},
		{
			desc:  "add invalid owner",
			owner: properties.Owner{},
			token: token,
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "update non-existant owner",
			owner: properties.Owner{Fname: "james", Lname: "Torredo", Phone: "0784677882"},
			token: token,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "update update owner with invalid token",
			owner: saved,
			token: wrongValue,
			err:   users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateOwner(tc.token, tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewOwner(t *testing.T) {
	svc := newService(map[string]string{token: email})

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	owner.ID, _ = svc.CreateOwner(token, owner)

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
			desc:     "view property with invalid token",
			identity: saved.ID,
			token:    wrongValue,
			err:      users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewOwner(tc.token, tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListOwners(t *testing.T) {
	svc := newService(map[string]string{token: email})

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.CreateOwner(token, owner)
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
			token:  token,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the properties",
			token:  token,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			token:  token,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			token:  token,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list properties with invalid token",
			token:  wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			err:    users.ErrUnauthorizedAccess,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListOwners(tc.token, tc.offset, tc.limit)
		size := uint64(len(page.Owners))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestFindOwner(t *testing.T) {
	svc := newService(map[string]string{token: email})

	fname := "james"
	lname := "torredo"
	phone := "0784677882"

	owner := properties.Owner{Fname: fname, Lname: lname, Phone: phone}

	_, err := svc.CreateOwner(token, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc  string
		token string
		fname string
		lname string
		phone string
		err   error
	}{
		{
			desc:  "find existing owner",
			token: token,
			fname: fname,
			lname: lname,
			phone: phone,
			err:   nil,
		},
		{
			desc:  "find owner with wrong first name",
			token: token,
			fname: "wrong",
			lname: lname,
			phone: phone,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with wrong last name",
			token: token,
			fname: fname,
			lname: "wrong",
			phone: phone,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with wrong phone number",
			token: token,
			fname: fname,
			lname: lname,
			phone: "wrong",
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with invalid token",
			token: wrongValue,
			fname: fname,
			lname: lname,
			phone: phone,
			err:   users.ErrUnauthorizedAccess,
		},
	}
	for _, tc := range cases {
		_, err := svc.FindOwner(tc.token, tc.fname, tc.lname, tc.phone)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
