package app

import (
	"context"
)

//AuthenticationService defines the authentication API.
type AuthenticationService interface {
	GetAuthToken(ctx *context.Context) error
	RenewAuthToken(ctx *context.Context) error
	RevokeAuthToken(ctx *context.Context) error
}
