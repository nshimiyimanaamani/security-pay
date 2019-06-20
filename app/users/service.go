package users

import (
	"errors"
	"github.com/rugwirobaker/paypack-backend/app/config"
	"github.com/rugwirobaker/paypack-backend/models"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	//ErrInvalidEntity indicates malformed entity specification (e.g.
	//invalid username,  password, account).
	ErrInvalidEntity = errors.New("invalid entity format specification")
)

//Service defines the users API
type Service interface {
	// Register creates new user account. In case of the failed registration, a
	// non-nil error value is returned.
	Register(models.Operator) (string, error)

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(models.Operator) (string, error)

	// Identify validates user's token. If token is valid, user's id
	// is returned. If token is invalid, or invocation failed for some
	// other reason, non-nil error values are returned in response.
	Identify(string) (string, error)
}

var _ Service = (*usersService)(nil)

type usersService struct {
	Config config.Config
	Idp    IdentityProvider
}

//New instanciates a new Service
func New(idp IdentityProvider, cfg config.Config) Service {
	return &usersService{Config: cfg, Idp: idp}
}

func (svc *usersService) Register(models.Operator) (string, error) { return "id", nil }

func (svc *usersService) Login(models.Operator) (string, error) { return "token", nil }

func (svc *usersService) Identify(string) (string, error) { return "id", nil }
