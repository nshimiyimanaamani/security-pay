package properties_test

import (
	"fmt"

	"testing"

	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/core/uuid"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	const op errors.Op = "app/properties/property.Validate"

	cases := []struct {
		desc     string
		property properties.Property
		err      error
	}{
		{
			desc: "validate with a valid property",
			property: properties.Property{
				Owner:      properties.Owner{ID: uuid.New().ID()},
				Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
				Due:        float64(1000),
				Namespace:  "kigali.gasabo.remera",
				RecordedBy: uuid.New().ID(),
			},
			err: nil,
		},
		{
			desc: "validate with empty owner field",
			property: properties.Property{
				Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
				Due:        float64(1000),
				Namespace:  "kigali.gasabo.remera",
				RecordedBy: uuid.New().ID(),
			},
			err: errors.E(op, "invalid property: missing owner", errors.KindBadRequest),
		},
		{
			desc: "validate with empty montly due",
			property: properties.Property{
				Owner:     properties.Owner{ID: uuid.New().ID()},
				Namespace: "kigali.gasabo.remera",
				Address:   properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
			},
			err: errors.E(op, "invalid property: missing due", errors.KindBadRequest),
		},
		{
			desc: "validate with empty address",
			property: properties.Property{
				Owner:      properties.Owner{ID: uuid.New().ID()},
				Address:    properties.Address{},
				Due:        float64(1000),
				Namespace:  "kigali.gasabo.remera",
				RecordedBy: uuid.New().ID(),
			},
			err: errors.E(op, "invalid property: invalid address", errors.KindBadRequest),
		},
		{
			desc: "validate with empty recorded by",
			property: properties.Property{
				Owner:   properties.Owner{ID: uuid.New().ID()},
				Address: properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
				Due:     float64(1000),
			},
			err: errors.E(op, "invalid property: missing recording agent", errors.KindBadRequest),
		},
		{
			desc: "validate with no account namespace",
			property: properties.Property{
				Owner:      properties.Owner{ID: uuid.New().ID()},
				Address:    properties.Address{Sector: "Remera", Cell: "Gishushu", Village: "Ingabo"},
				Due:        float64(1000),
				RecordedBy: uuid.New().ID(),
			},
			err: errors.E(op, "invalid property: missing namespace tag", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.property.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
