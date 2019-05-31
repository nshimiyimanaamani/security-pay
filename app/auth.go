package app

import (
	"context"
)

//AuthenticationService defines the authentication API.
type AuthenticationService interface {
	GetAuthToken(ctx *context.Context, r *GetAuthTokenRequest) error
	RenewAuthToken(ctx *context.Context, r *RenewAuthTokenRequest) error
	RevokeAuthToken(ctx *context.Context, r *RevokeAuthTokenRequest) error
}

//GetAuthTokenRequest defines a request to the GetAuthToken Endpoint.
type GetAuthTokenRequest struct{}

//RenewAuthTokenRequest defines a request to the RenewAuthToken Endpoint.
type RenewAuthTokenRequest struct{}

//RevokeAuthTokenRequest defines a request to the RevokeAuthToken Endpoint.
type RevokeAuthTokenRequest struct{}
