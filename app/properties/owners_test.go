package properties_test

import (
	"fmt"

	"testing"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/stretchr/testify/assert"
)

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
