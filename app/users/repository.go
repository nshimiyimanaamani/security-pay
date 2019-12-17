package users

import "context"

// Repository is the users database interface
type Repository interface {
	AdminsRepository
	AgentsRepository
	DevelopersRepository
	ManagersRepository
}

// AdminsRepository ...
type AdminsRepository interface {
	SaveAdmin(ctx context.Context, user Administrator) (Administrator, error)
	RetrieveAdmin(ctx context.Context, id string) (Administrator, error)
	ListAdmins(ctx context.Context, offset, limit uint64) (AdministratorPage, error)
	UpdateAdminCreds(ctx context.Context, user Administrator) error
}

// AgentsRepository ...
type AgentsRepository interface {
	SaveAgent(ctx context.Context, user Agent) (Agent, error)
	RetrieveAgent(ctx context.Context, id string) (Agent, error)
	ListAgents(ctx context.Context, offset, limit uint64) (AgentPage, error)
	UpdateAgentDetails(ctx context.Context, user Agent) error
	UpdateAgentCreds(ctx context.Context, user Agent) error
}

// DevelopersRepository ...
type DevelopersRepository interface {
	SaveDeveloper(ctx context.Context, user Developer) (Developer, error)
	RetrieveDeveloper(ctx context.Context, id string) (Developer, error)
	ListDevelopers(ctx context.Context, offset, limit uint64) (DeveloperPage, error)
	UpdateDeveloperCreds(ctx context.Context, user Developer) error
}

// ManagersRepository ...
type ManagersRepository interface {
	SaveManager(ctx context.Context, user Manager) (Manager, error)
	RetrieveManager(ctx context.Context, id string) (Manager, error)
	ListManagers(ctx context.Context, offset, limit uint64) (ManagerPage, error)
	UpdateManagerCreds(ctx context.Context, user Manager) error
}
