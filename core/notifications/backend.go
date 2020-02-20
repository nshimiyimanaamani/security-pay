package notifications

import "context"

// Backend abstracts away the sms client
type Backend interface {
	Send(ctx context.Context, id, message string, recipients []string) error
}
