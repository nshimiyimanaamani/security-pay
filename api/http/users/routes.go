package users

// admin routes
const (
	RegisterAdminRoute    = "/accounts/admin"
	RetrieveAdminRoute    = "/accounts/admin/{email}"
	UpdateAdminCredsRoute = "/accounts/{account}/admin{email}"
)

// agents routes
const (
	RegisterAgentRoute    = "/accounts/agents"
	RetrieveAgentRoute    = "/accounts/admin/{phone}"
	UpdateAgentRoute      = "/accounts/agents/{phone}"
	UpdateAgentCredsRoute = "/accounts/agents/{phone}"
	ListAgentsRoute       = "/accounts/agents"
)

// developer routes
const (
	RegisterDeveloperRoute    = "/accounts/developers"
	RetrieveDeveloperRoute    = "/accounts/developers/{email}"
	UpdateDeveloperCredsRoute = "/accounts/developers/{email}"
	ListDevelopersRoute       = "/accounts/developers"
)

// namnager routes
const (
	RegisterManagerRoute    = "/accounts/managers"
	RetrieveManagerRoute    = "/accounts/managers/{email}"
	UpdateManagerCredsRoute = "/accounts/managers{email}"
	ListManagersRoute       = "/accounts/managers"
)
