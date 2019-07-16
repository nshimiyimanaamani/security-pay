package properties_test

import (
	"fmt"

	"testing"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
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
		{"validate with a valid property", property, nil},
		{"validate with empty owner field", invalidProperty, properties.ErrInvalidEntity},
		{"validate with empty montly due", emptyDue, properties.ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.property.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
