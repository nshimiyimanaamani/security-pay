package sms_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/app/sms"
)
var message = models.Message{ Body:"hello world", Destination:"+250-788-455-100"}

func newService()sms.Service{
	return sms.New()
}

func TestSend(t *testing.T){
	svc:= newService()

	cases:= []struct{
		desc 		string
		message 	models.Message
		err 		error
	}{
		{"send valid message", message, nil},
		{"send invalid message", models.Message{Body:"", Destination:"+250-788-455-100"}, models.ErrInvalidEntity},
	}

	for _,tc:=range cases{
		err:=svc.Send(tc.message)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}