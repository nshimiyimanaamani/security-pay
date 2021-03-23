package users

// admin routes
const (
	RegisterAdminRoute    = "/accounts/admin"
	ListAdminsRoute       = "/accounts/admins"
	RetrieveAdminRoute    = "/accounts/admin/{email}"
	UpdateAdminCredsRoute = "/accounts/admin/creds/{email}"
)

// agents routes
const (
	RegisterAgentRoute    = "/accounts/agents"
	RetrieveAgentRoute    = "/accounts/agents/{phone}"
	UpdateAgentRoute      = "/accounts/agents/{phone}"
	DeleteAgentRoute      = "/accounts/agents/{phone}"
	UpdateAgentCredsRoute = "/accounts/agents/creds/{phone}"
	ListAgentsRoute       = "/accounts/agents"
)

// developer routes
const (
	RegisterDeveloperRoute    = "/accounts/developers"
	RetrieveDeveloperRoute    = "/accounts/developers/{email}"
	DeleteDeveloperRoute      = "/accounts/developers/{email}"
	UpdateDeveloperCredsRoute = "/accounts/developers/creds/{email}"
	ListDevelopersRoute       = "/accounts/developers"
)

// namnager routes
const (
	RegisterManagerRoute    = "/accounts/managers"
	RetrieveManagerRoute    = "/accounts/managers/{email}"
	DeleteManagerRoute      = "/accounts/managers/{email}"
	UpdateManagerCredsRoute = "/accounts/managers/creds/{email}"
	ListManagersRoute       = "/accounts/managers"
)
