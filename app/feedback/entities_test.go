package feedback

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	const op errors.Op = "app/feedback/message.Validate"

	cases := []struct {
		desc string
		msg  Message
		err  error
	}{
		{
			desc: "valid message",
			msg:  Message{Title: "title", Body: "body", Creator: "0784677882"},
			err:  nil,
		},
		{
			desc: "message with no title",
			msg:  Message{Body: "body", Creator: "0784677882"},
			err:  errors.E(op, "invalid message: missing title", errors.KindBadRequest),
		},
		{
			desc: "message with no body",
			msg:  Message{Title: "title", Creator: "0784677882"},
			err:  errors.E(op, "invalid message: missing body", errors.KindBadRequest),
		},
		{
			desc: "message with invalid creator(phone number)",
			msg:  Message{Title: "title", Body: "body", Creator: "8888"},
			err:  errors.E(op, "invalid message: invalid phone number", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.msg.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
