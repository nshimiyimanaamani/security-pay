package models

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	id		 = "1000-4433-34343"
	amount   = "5000.00"
	method   = "BK" 
	property = "1000-4433-34343"
)

func TestTransactionValidate(t *testing.T){
	cases:= []struct{
		desc string
		transaction Transaction
		err  error
	}{
		{"validate user with valid data", Transaction{ID:id, Amount:amount, Method:method, Property:property}, nil},
		{"validate user with invalid data", Transaction{}, ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.transaction.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}