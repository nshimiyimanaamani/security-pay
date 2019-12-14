package user

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Developer represents users of type developer
type Developer struct {
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Account   string    `json:"account,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate developer
func (dev *Developer) Validate() error {
	const op errors.Op = "app/accounts/developer.Validate"

	if dev.Email == "" {
		return errors.E(op, "invalid developer: missing email", errors.KindBadRequest)
	}

	if dev.Password == "" {
		return errors.E(op, "invalid developer: missing password", errors.KindBadRequest)
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
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate agent
func (ag *Agent) Validate() error {
	const op errors.Op = "app/accounts/agent.Validate"

	if ag.Telephone == "" {
		return errors.E(op, "invalid agent: missing email", errors.KindBadRequest)
	}
	if ag.FirstName == "" || ag.LastName == "" {
		return errors.E(op, "invalid agent: missing names", errors.KindBadRequest)
	}
	if ag.Cell == "" {
		return errors.E(op, "invalid agent: missing cell", errors.KindBadRequest)
	}
	if ag.Sector == "" {
		return errors.E(op, "invalid agent: missing sector", errors.KindBadRequest)
	}
	if ag.Village == "" {
		return errors.E(op, "invalid agent: missing village", errors.KindBadRequest)
	}
	if ag.Account == "" {
		return errors.E(op, "invalid agent: missing account reference", errors.KindBadRequest)
	}
	return nil
}

// Manager represent users of type manager
type Manager struct {
	Email     string    `json:"email,omitempty"`
	Cell      string    `json:"cell,omitempty"`
	Password  string    `json:"password,omitempty"`
	Account   string    `json:"account,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate manager
func (m *Manager) Validate() error {
	const op errors.Op = "app/accounts/manager.Validate"
	if m.Email == "" {
		return errors.E(op, "invalid manager: missing email", errors.KindBadRequest)
	}
	if m.Cell == "" {
		return errors.E(op, "invalid manager: missing email", errors.KindBadRequest)
	}
	if m.Account == "" {
		return errors.E(op, "invalid manager: missing account reference", errors.KindBadRequest)
	}
	return nil
}

// DeveloperPage is a collection of developers and metadata.
type DeveloperPage struct {
	Developers []Developer
	PageMetadata
}

// AgentPage is a collection of agents plus metadata.
type AgentPage struct {
	Agents []Agent
	PageMetadata
}

// ManagerPage is a collection of managers plus metadata.
type ManagerPage struct {
	Managers []Manager
	PageMetadata
}
