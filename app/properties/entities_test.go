package properties_test

import (
	"fmt"

	"testing"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
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
		err      error
	}{
		{"validate with a valid property", property, nil},
		{"validate with empty owner field", invalidProperty, properties.ErrInvalidEntity},
		{"validate with empty montly due", emptyDue, properties.ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.property.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestValidateOwner(t *testing.T) {
	cases := []struct {
		desc  string
		owner properties.Owner
		err   error
	}{
		{
			desc:  "validate with a valid owner entity",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"},
			err:   nil,
		},
		{
			desc:  "validate with empty fname field",
			owner: properties.Owner{Lname: "Torredo", Phone: "0784677882"},
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "validate with empty lname field",
			owner: properties.Owner{Fname: "James", Phone: "0784677882"},
			err:   properties.ErrInvalidEntity,
		},
		{
			desc:  "validate with invalid phone number",
			owner: properties.Owner{Fname: "James", Lname: "Torredo", Phone: "77878333"},
			err:   properties.ErrInvalidEntity,
		},
	}

	for _, tc := range cases {
		err := tc.owner.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
