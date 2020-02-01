package accounts

import "context"

// Repository is the account database adapter
type Repository interface {
	Save(ctx context.Context, acc Account) (Account, error)
	Update(ctx context.Context, acc Account) error
	Retrieve(ctx context.Context, id string) (Account, error)
	List(ctx context.Context, offset, limit uint64) (AccountPage, error)
}
