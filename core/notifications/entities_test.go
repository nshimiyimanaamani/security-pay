package notifications_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/notifications"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidatePayload(t *testing.T) {
	const op errors.Op = "core/notifications/Payload.Validate"

	cases := []struct {
		desc    string
		payload notifications.Payload
		err     error
	}{
		{
			desc:    "validate valid account",
			payload: notifications.Payload{Message: "message", Recipients: []string{"recipient"}},
			err:     nil,
		},
		{
			desc:    "validate valid account",
			payload: notifications.Payload{Recipients: []string{"recipient"}},
			err:     errors.E(op, "invalid payload: message is required", errors.KindBadRequest),
		},
		{
			desc:    "validate valid account",
			payload: notifications.Payload{Message: "message"},
			err:     errors.E(op, "invalid payload: recipients must be a non empty array", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.payload.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
