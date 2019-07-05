package properties_test

import (
	"fmt"

	"testing"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/stretchr/testify/assert"
)
func TestValidate(t *testing.T){
	cases:= []struct{
		desc string
		property properties.Property
		err  error
	}{
		{"validate with a valid property", property, nil},
		{"validate with empty owner field", wrongProperty, properties.ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.property.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}