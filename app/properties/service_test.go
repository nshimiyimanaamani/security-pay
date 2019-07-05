package properties_test

import (
	"fmt"
	"testing"

	//"github.com/rugwirobaker/paypack-backend/app"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	property = properties.Property{
		Owner: "Eugene Mugabo",
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
	}

	wrongProperty = properties.Property{
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
	}

	wrongValue = "wrong-value"
)

func newService() properties.Service {
	idp := mocks.NewIdentityProvider()
	store := mocks.NewPropertyStore()
	return properties.New(idp, store)
}

func TestAddProperty(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{"add valid property", property, nil},
		{"add invalid property", wrongProperty, properties.ErrInvalidEntity},
	}

	for _, tc := range cases {
		_, err := svc.AddProperty(tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService()

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
			property: wrongProperty,
			err:      properties.ErrInvalidEntity,
		},
		{
			desc:     "update non-existant property",
			property: property,
			err:      properties.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := svc.UpdateProperty(tc.property)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewProperty(t *testing.T) {
	svc := newService()
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

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		svc.AddProperty(property)
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
