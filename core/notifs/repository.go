package notifs

import "context"

// Repository commits sent messages to the database
type Repository interface {
	// Saves sent message info to the database
	Save(ctx context.Context, sms Notification) (Notification, error)

	// Find a message by id
	Find(ctx context.Context, id string) (Notification, error)

	// List messages by namespace(account)
	List(ctx context.Context, namespace string, offset, limit uint64) (NoticationPage, error)

	// Count messages by namespace(account)
	Count(ctx context.Context, namespace string) (uint64, error)
}
