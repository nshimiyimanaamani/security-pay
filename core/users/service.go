package users

import (
	"context"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/encrypt"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/passwords"
)

// Service defines users usecases
type Service interface {
	Admins
	Agents
	Developers
	Managers
}

// Admins ...
type Admins interface {
	DeleteAdmin(ctx context.Context, id string) error
	RegisterAdmin(ctx context.Context, user Administrator) (Administrator, error)
	RetrieveAdmin(ctx context.Context, id string) (Administrator, error)
	ListAdmins(ctx context.Context, offset, limit uint64) (AdministratorPage, error)
	UpdateAdminCreds(ctx context.Context, user Administrator) error
}

// Agents ...
type Agents interface {
	DeleteAgent(ctx context.Context, id string) error
	RegisterAgent(ctx context.Context, user Agent) (Agent, error)
	RetrieveAgent(ctx context.Context, id string) (Agent, error)
	ListAgents(ctx context.Context, offset, limit uint64) (AgentPage, error)
	UpdateAgent(ctx context.Context, user Agent) error
	UpdateAgentCreds(ctx context.Context, user Agent) error
}

// Developers ...
type Developers interface {
	DeleteDeveloper(ctx context.Context, id string) error
	RegisterDeveloper(ctx context.Context, user Developer) (Developer, error)
	RetrieveDeveloper(ctx context.Context, id string) (Developer, error)
	ListDevelopers(ctx context.Context, offset, limit uint64) (DeveloperPage, error)
	UpdateDeveloperCreds(ctx context.Context, user Developer) error
}

// Managers ...
type Managers interface {
	DeleteManager(ctx context.Context, id string) error
	RegisterManager(ctx context.Context, user Manager) (Manager, error)
	RetrieveManager(ctx context.Context, id string) (Manager, error)
	ListManagers(ctx context.Context, offset, limit uint64) (ManagerPage, error)
	UpdateManagerCreds(ctx context.Context, user Manager) error
}

type service struct {
	pgen      passwords.Generator
	hasher    passwords.Hasher
	encrypter encrypt.Encrypter
	repo      Repository
}

// Options ...
type Options struct {
	PGen      passwords.Generator
	Hasher    passwords.Hasher
	Encrypter encrypt.Encrypter
	Repo      Repository
}

// New creates an instance of users.Service
func New(opts *Options) Service {
	return &service{
		hasher:    opts.Hasher,
		encrypter: opts.Encrypter,
		repo:      opts.Repo,
		pgen:      opts.PGen,
	}
}
