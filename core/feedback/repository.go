package feedback

import "context"

// Repository is the API to the feedback database(table)
type Repository interface {
	Save(ctx context.Context, msg *Message) (*Message, error)
	Retrieve(ctx context.Context, id string) (Message, error)
	RetrieveAll(ctx context.Context, offset, limit uint64) (MessagePage, error)
	Update(ctx context.Context, msg Message) error
	Delete(ctx context.Context, id string) error
}
