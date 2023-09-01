package users

import (
	"time"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// Role represents user access level
type Role string

const (
	// Dev has access oveer all accounts
	Dev = "dev"
	// Admin has access level only to the sector they manage
	Admin = "admin"
	// Basic has access to only a single cell.
	Basic = "basic"

	// Min is the minimun privilage level
	Min = "min"
)

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Administrator ...
type Administrator struct {
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Account   string    `json:"account,omitempty"`
	Role      Role      `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate ent
func (ent *Administrator) Validate() error {
	const op errors.Op = "app/accounts/administrator.Validate"

	if ent.Email == "" {
		return errors.E(op, "invalid manager: missing email", errors.KindBadRequest)
	}
	if ent.Account == "" {
		return errors.E(op, "invalid manager: missing account", errors.KindBadRequest)
	}
	return nil
}

// Agent represents users of type agents
type Agent struct {
	Telephone string    `json:"telephone,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Cell      string    `json:"cell,omitempty"`
	Sector    string    `json:"sector,omitempty"`
	Village   string    `json:"village,omitemty"`
	Password  string    `json:"password,omitempty"`
	Account   string    `json:"account,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate agent
func (ent *Agent) Validate() error {
	const op errors.Op = "app/accounts/agent.Validate"

	if ent.Telephone == "" {
		return errors.E(op, "invalid agent: missing phone", errors.KindBadRequest)
	}
	if ent.FirstName == "" || ent.LastName == "" {
		return errors.E(op, "invalid agent: missing names", errors.KindBadRequest)
	}
	if ent.Cell == "" || ent.Sector == "" || ent.Village == "" {
		return errors.E(op, "invalid agent: missing address references", errors.KindBadRequest)
	}
	if ent.Account == "" {
		return errors.E(op, "invalid manager: missing account", errors.KindBadRequest)
	}
	return nil
}

// Developer ...
type Developer struct {
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Account   string    `json:"account,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate developer
func (ent *Developer) Validate() error {
	const op errors.Op = "app/accounts/developer.Validate"

	if ent.Email == "" {
		return errors.E(op, "invalid manager: missing email", errors.KindBadRequest)
	}
	if ent.Password == "" {
		return errors.E(op, "invalid manager: missing password", errors.KindBadRequest)
	}
	if ent.Account == "" {
		return errors.E(op, "invalid manager: missing account", errors.KindBadRequest)
	}
	return nil
}

// Manager represent users of type manager
type Manager struct {
	Email     string    `json:"email,omitempty"`
	Cell      string    `json:"cell,omitempty"`
	Password  string    `json:"password,omitempty"`
	Account   string    `json:"account,omitempty"`
	Role      string    `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate manager
func (ent *Manager) Validate() error {
	const op errors.Op = "app/accounts/manager.Validate"
	if ent.Email == "" {
		return errors.E(op, "invalid manager: missing email", errors.KindBadRequest)
	}
	if ent.Cell == "" {
		return errors.E(op, "invalid manager: missing cell", errors.KindBadRequest)
	}
	if ent.Account == "" {
		return errors.E(op, "invalid manager: missing account", errors.KindBadRequest)
	}
	return nil
}

// AdministratorPage is a collection of administrators + metadata
type AdministratorPage struct {
	Administrators []Administrator
	PageMetadata
}

// AgentPage is a collection of agents + metadata.
type AgentPage struct {
	Agents []Agent
	PageMetadata
}

// DeveloperPage is a collection of developers + metadata.
type DeveloperPage struct {
	Developers []Developer
	PageMetadata
}

// ManagerPage is a collection of managers + metadata.
type ManagerPage struct {
	Managers []Manager
	PageMetadata
}
