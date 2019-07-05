package sms

import "errors"

var (
	// ErrConflict attempt to create an entity with an alreasdy existing id
	ErrConflict = errors.New("message already exists")
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid message format")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent message")
)

//Service defines the send sms API
type Service interface {
	//Send sends a message and returns nil if the operation is successful
	Send(Message) error
}

var _ (Service) = (*smsService)(nil)

type smsService struct{}

//New instanciates a new Service instance
func New() Service {
	return &smsService{}
}

func (svc *smsService) Send(message Message) error {
	return message.Validate()
}
