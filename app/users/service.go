package users

import (
	"errors"

	"github.com/rugwirobaker/paypack-backend/app"
)

var (
	// ErrConflict attempt to create an entity with an alreasdy existing id
	ErrConflict = errors.New("entity already exists")
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid entity format")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
)

//Service defines the users API
type Service interface {
	// Register creates new user account. In case of the failed registration, a
	// non-nil error value is returned.
	Register(User) (string, error)

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(User) (string, error)

	// Identify validates user's token. If token is valid, user's id
	// is returned. If token is invalid, or invocation failed for some
	// other reason, non-nil error values are returned in response.
	Identify(string) (string, error)
}

var _ Service = (*usersService)(nil)

type usersService struct {
	hasher  Hasher
	tempIdp TempIdentityProvider
	idp     app.IdentityProvider
	store   Store
}

//New instanciates a new Service.
func New(hasher Hasher, tempIdp TempIdentityProvider, idp app.IdentityProvider, store Store) Service {
	return &usersService{hasher: hasher, tempIdp: tempIdp, idp: idp, store: store}
}

func (svc *usersService) Register(user User) (string, error) {
	hash, err := svc.hasher.Hash(user.Password)
	if err != nil {
		return "", ErrInvalidEntity
	}

	if !user.CheckCell() {
		return "", ErrInvalidEntity
	}

	user.Password = hash

	user.ID = svc.idp.ID()

	return svc.store.Save(user)
}

func (svc *usersService) Login(user User) (string, error) {
	dbUser, err := svc.store.RetrieveByID(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := svc.hasher.Compare(user.Password, dbUser.Password); err != nil {
		return "", ErrUnauthorizedAccess
	}

	return svc.tempIdp.TemporaryKey(user.Email)
}

func (svc *usersService) Identify(token string) (string, error) {
	id, err := svc.tempIdp.Identity(token)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	return id, nil
}
