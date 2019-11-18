package feedback

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		desc string
		msg  Message
		err  error
	}{
		{"valid message", Message{Title: "title", Body: "body"}, nil},
		{"message with no title", Message{Body: "body"}, ErrInvalidEntity},
		{"message with no body", Message{Title: "title"}, ErrInvalidEntity},
	}

	for _, tc := range cases {
		err := tc.msg.Validate()
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
