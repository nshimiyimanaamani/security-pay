package sms_test


import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/sms"
	"github.com/stretchr/testify/assert"
)

var (
	body = "hello world"
	destination = "+250-788-455-100"
)
func TestValidate(t *testing.T){
	cases:= []struct{
		desc string
		message sms.Message
		err  error
	}{
		{"validate with valid data", sms.Message{Body:body, Destination: destination}, nil},
		{"validate with empty message", sms.Message{}, sms.ErrInvalidEntity},
		//TODO add more validation cases
	}

	for _, tc := range cases {
		err := tc.message.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}

}