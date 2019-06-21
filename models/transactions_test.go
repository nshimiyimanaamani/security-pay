package models

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	amount    = "5000.00"
	property = "1000-4433-34343"
)

func TestTransactionValidate(t *testing.T){
	cases:= []struct{
		desc string
		transaction Transaction
		err  error
	}{
		{"validate user with valid data", Transaction{Amount:amount,Property:property}, nil},
	}

	for _, tc := range cases {
		err := tc.transaction.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}