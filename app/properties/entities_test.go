package properties_test

import (
	"fmt"

	"testing"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{
			desc: "validate with a valid property",
			property: properties.Property{
				Owner:   properties.Owner{ID: uuid.New().ID()},
				Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
				Due:     float64(1000),
			},
			err: nil,
		},
		{
			desc: "validate with empty owner field",
			property: properties.Property{
				Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
				Due:     float64(1000),
			},
			err: properties.ErrInvalidEntity,
		},
		{
			desc: "validate with empty montly due",
			property: properties.Property{
				Owner:   properties.Owner{ID: uuid.New().ID()},
				Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
			},
			err: properties.ErrInvalidEntity,
		},
		{
			desc: "validate with empty address",
			property: properties.Property{
				Owner:   properties.Owner{ID: uuid.New().ID()},
				Address: properties.Address{},
				Due:     float64(1000),
			},
			err: properties.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		err := tc.property.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
