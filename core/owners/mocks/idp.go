package mocks

import (
	"fmt"
	"sync"

	"github.com/rugwirobaker/paypack-backend/core/identity"
)

var _ identity.Provider = (*identityProviderMock)(nil)

type identityProviderMock struct {
	mu      sync.Mutex
	counter int
}

// NewIdentityProvider creates "mirror" identity provider, i.e. generated
// token will hold value provided by the caller.
func NewIdentityProvider() identity.Provider {
	return &identityProviderMock{}
}

func (idp *identityProviderMock) ID() string {
	idp.mu.Lock()
	defer idp.mu.Unlock()

	idp.counter++
	return fmt.Sprintf("%s%012d", "123e4567-e89b-12d3-a456-", idp.counter)
}
