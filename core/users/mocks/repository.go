package mocks

import (
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/users"
)

var _ users.Repository = (*userRepository)(nil)

type userRepository struct {
	mu        sync.Mutex
	adcounter uint64
	agcounter uint64
	dcounter  uint64
	mcounter  uint64
	accounts  map[string]string
	admins    map[string]users.Administrator
	agents    map[string]users.Agent
	devs      map[string]users.Developer
	managers  map[string]users.Manager
}

// NewRepository ...
func NewRepository(accounts map[string]string) users.Repository {
	return &userRepository{
		accounts: accounts,
		admins:   make(map[string]users.Administrator),
		agents:   make(map[string]users.Agent),
		devs:     make(map[string]users.Developer),
		managers: make(map[string]users.Manager),
	}
}
