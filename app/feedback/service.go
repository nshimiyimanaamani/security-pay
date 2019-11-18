package feedback

import (
	"context"
	"errors"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/identity"
)

var (
	// ErrConflict attempt to create an entity with an alreasdy existing id
	ErrConflict = errors.New("message already exists")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid entity format")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
)

// Service ...
type Service interface {
	Record(ctx context.Context, msg Message) error
	Retrieve(ctx context.Context, id string) (Message, error)
	Update(ctx context.Context, msg Message) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	idp  identity.Provider
	repo Repository
}

func (svc *service) Record(ctx context.Context, msg Message) error {

	msg.Validate()

	msg = *svc.newMsg(&msg)

	return svc.repo.Save(ctx, msg)
}
func (svc *service) Retrieve(ctx context.Context, id string) (Message, error) {
	return svc.repo.Retrieve(ctx, id)
}
func (svc *service) Update(ctx context.Context, msg Message) error {
	return svc.repo.Update(ctx, msg)
}
func (svc *service) Delete(ctx context.Context, id string) error {
	return svc.repo.Delete(ctx, id)
}

func (svc *service) newMsg(msg *Message) *Message {
	var timestamp = time.Now()

	msg.ID = svc.idp.ID()

	msg.CreatedAt = timestamp
	msg.UpdatedAt = timestamp

	return msg
}
