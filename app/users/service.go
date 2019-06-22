package users

import (
	"github.com/rugwirobaker/paypack-backend/app/config"
	"github.com/rugwirobaker/paypack-backend/models"
	store "github.com/rugwirobaker/paypack-backend/store/users"
)

//Service defines the users API
type Service interface {
	// Register creates new user account. In case of the failed registration, a
	// non-nil error value is returned.
	Register(models.User) (string, error)

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(models.User) (string, error)

	// Identify validates user's token. If token is valid, user's id
	// is returned. If token is invalid, or invocation failed for some
	// other reason, non-nil error values are returned in response.
	Identify(string) (string, error)
}

var _ Service = (*usersService)(nil)

type usersService struct {
	Config config.Config
	hasher Hasher
	idp    IdentityProvider
	store  store.UserStore
}

//New instanciates a new Service.
func New(idp IdentityProvider, cfg config.Config, hasher Hasher, store store.UserStore) Service {
	return &usersService{Config: cfg, hasher: hasher, idp: idp, store: store}
}

func (svc *usersService) Register(user models.User) (string, error) {
	hash, err := svc.hasher.Hash(user.Password)
	if err != nil {
		return "", models.ErrInvalidEntity
	}
	user.Password = hash
	return svc.store.Save(user)
}

func (svc *usersService) Login(user models.User) (string, error) {
	dbUser, err := svc.store.RetrieveByID(user.Email)
	if err != nil {
		return "", models.ErrUnauthorizedAccess
	}

	if err := svc.hasher.Compare(user.Password, dbUser.Password); err != nil {
		return "", models.ErrUnauthorizedAccess
	}

	return svc.idp.TemporaryKey(user.Email)
}

func (svc *usersService) Identify(token string) (string, error) {
	id, err := svc.idp.Identity(token)
	if err != nil {
		return "", models.ErrUnauthorizedAccess
	}
	return id, nil
}
