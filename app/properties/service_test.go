package properties_test

import (
	"fmt"
	"testing"

	//"github.com/rugwirobaker/paypack-backend/app"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	wrongValue = "wrong-value"
)

func newService() properties.Service {
	idp := mocks.NewIdentityProvider()
	propStore := mocks.NewPropertyStore()
	ownerStore := mocks.NewOwnerStore()
	return properties.New(idp, ownerStore, propStore)
}

func TestAddProperty(t *testing.T) {
	svc := newService()

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
		err      error
	}{
		{"add valid property", property, nil},
		{"add invalid property", invalidProperty, properties.ErrInvalidEntity},
		{"validate with empty montly due", emptyDue, properties.ErrInvalidEntity},
	}

	for _, tc := range cases {
		_, err := svc.AddProperty(tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService()

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

	saved, _ := svc.AddProperty(property)

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
			err:      properties.ErrNotFound,
		},
		{
			desc:     "update property with empty due",
			property: emptyDue,
			err:      properties.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateProperty(tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewProperty(t *testing.T) {
	svc := newService()

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	saved, _ := svc.AddProperty(property)

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
			err:      properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewProperty(tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByOwner(t *testing.T) {
	svc := newService()

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}
	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}

	oid, err := svc.CreateOwner(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		property.Owner = oid
		_, err = svc.AddProperty(property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := []struct {
		desc   string
		owner  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all properties",
			owner:  property.Owner,
			offset: 0,
			limit:  n,
			size:   n,
		},
		{
			desc:   "list half of the properties",
			owner:  property.Owner,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			owner:  property.Owner,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			owner:  property.Owner,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesByOwner(tc.owner, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesBySector(t *testing.T) {
	svc := newService()

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(property)
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
			sector: property.Sector,
			offset: 0,
			limit:  n,
			size:   n,
		},
		{
			desc:   "list half of the properties",
			sector: property.Sector,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			sector: property.Sector,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			sector: property.Sector,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesBySector(tc.sector, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListPropertiesByCell(t *testing.T) {
	svc := newService()

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(property)
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
			cell:   property.Cell,
			offset: 0,
			limit:  n,
			size:   n,
		},
		{
			desc:   "list half of the properties",
			cell:   property.Cell,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			cell:   property.Cell,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			cell:   property.Cell,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesByCell(tc.cell, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestListPropertiesByVillage(t *testing.T) {
	svc := newService()

	property := properties.Property{
		Owner: "Eugene Mugabo", Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
		Due: float64(1000),
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(property)
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
			village: property.Village,
			offset:  0,
			limit:   n,
			size:    n,
		},
		{
			desc:    "list half of the properties",
			village: property.Village,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
			err:     nil,
		},
		{
			desc: "	list empty set",
			village: property.Village,
			offset:  n + 1,
			limit:   n,
			size:    0,
			err:     nil,
		},
		{
			desc:    "list with zero limit",
			village: property.Village,
			offset:  1,
			limit:   0,
			size:    0,
			err:     nil,
		},
	}

	for _, tc := range cases {
		page, err := svc.ListPropertiesByVillage(tc.village, tc.offset, tc.limit)
		size := uint64(len(page.Properties))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}

}

func TestCreateOwner(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc  string
		owner properties.Owner
		err   error
	}{
		{
			desc:  "add valid owner",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"},
			err:   nil,
		},
		{
			desc:  "add invalid owner",
			owner: properties.Owner{},
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with empty fname field",
			owner: properties.Owner{Lname: "Torredo", Phone: "0784677882"},
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with empty lname field",
			owner: properties.Owner{Fname: "James", Phone: "0784677882"},
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "add owner with invalid phone number",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "77878333"},
			err:   properties.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		_, err := svc.CreateOwner(tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateOwner(t *testing.T) {
	svc := newService()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	owner.ID, _ = svc.CreateOwner(owner)

	saved := owner

	cases := []struct {
		desc  string
		owner properties.Owner
		err   error
	}{
		{
			desc:  "update existing owner",
			owner: saved,
			err:   nil,
		},
		{
			desc:  "update non-existant owner",
			owner: properties.Owner{Fname: "james", Lname: "Torredo", Phone: "0784677882"},
			err:   properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateOwner(tc.owner)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewOwner(t *testing.T) {
	svc := newService()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	owner.ID, _ = svc.CreateOwner(owner)

	saved := owner

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
			err:      properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewOwner(tc.identity)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestListOwners(t *testing.T) {
	svc := newService()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.CreateOwner(owner)
	}

	cases := []struct {
		desc   string
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
		page, err := svc.ListOwners(tc.offset, tc.limit)
		size := uint64(len(page.Owners))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestFindOwner(t *testing.T) {
	svc := newService()

	fname := "james"
	lname := "torredo"
	phone := "0784677882"

	owner := properties.Owner{Fname: fname, Lname: lname, Phone: phone}

	_, err := svc.CreateOwner(owner)
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
			fname: fname,
			lname: lname,
			phone: phone,
			err:   nil,
		},
		{
			desc:  "find owner with wrong first name",
			fname: "wrong",
			lname: lname,
			phone: phone,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with wrong last name",
			fname: fname,
			lname: "wrong",
			phone: phone,
			err:   properties.ErrNotFound,
		},
		{
			desc:  "find owner with wrong phone number",
			fname: fname,
			lname: lname,
			phone: "wrong",
			err:   properties.ErrNotFound,
		},
	}
	for _, tc := range cases {
		_, err := svc.FindOwner(tc.fname, tc.lname, tc.phone)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
