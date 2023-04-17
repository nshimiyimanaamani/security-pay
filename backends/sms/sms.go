package sms

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/sms-client/fdi"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Timeout sets the default client timeout
const Timeout = 30 * time.Second

var _ (notifs.Backend) = (*Backend)(nil)

// Backend ..
type Backend struct {
	URL      string
	Token    string
	Refresh  string
	SenderID string
	client   *fdi.Client
}

// Options ...
type Options struct {
	URL       string
	AppID     string
	AppSecret string
	SenderID  string
}

func New(opts *Options) (*Backend, error) {
	client, err := fdi.NewDefault(opts.AppID, opts.AppSecret)
	if err != nil {
		return nil, fmt.Errorf("fdi: client sms failed to inialized: %w", err)
	}
	return &Backend{
		URL:      opts.URL,
		client:   client,
		SenderID: "PayPack",
	}, nil
}

// Send sends sms to the specified number
func (s *Backend) Send(ctx context.Context, id, message string, recipients []string) error {
	const op errors.Op = "backends/sms/Send"

	in := fdi.Message{
		ID:         id,
		Body:       message,
		Recipients: recipients,
		Sender:     s.SenderID,
		Report:     s.URL,
	}

	_, _, err := s.client.Send(ctx, in)
	if err != nil {
		return errors.E(op, err, nil)
	}
	return nil
}
