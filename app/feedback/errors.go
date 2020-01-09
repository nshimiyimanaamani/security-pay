package feedback

import "errors"

var (
	// ErrConflict attempt to create an entity with an alreasdy existing id
	ErrConflict = errors.New("message already exists")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid message entity")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("message entity not found")
)
