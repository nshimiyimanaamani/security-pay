package mocks

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type mockSmsBackend struct{}

// NewBackend creates a new notifs.Backend mock
func NewBackend() notifs.Backend {
	return &mockSmsBackend{}
}

func (repo *mockSmsBackend) Send(ctx context.Context, id, message string, recipients []string) error {
	const op errors.Op = "core/notifs/mocks/mockSmsBackend.Send"

	if recipients == nil || len(recipients) < 1 {
		return errors.E(op, "the number of recipients must be greater than zero", errors.KindBadRequest)
	}

	if id == "" {
		return errors.E(op, "invalid sms id", errors.KindBadRequest)
	}

	if message == "" {
		return errors.E(op, "sms body must not be empty", errors.KindBadRequest)
	}
	return nil
}
