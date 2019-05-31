package app

import (
	"context"
)

//AuthenticationService defines the authentication API.
type AuthenticationService interface {
	GetAuthToken(ctx *context.Context, r *GetAuthTokenRequest) (GetAuthTokenResponse, error)
	RenewAuthToken(ctx *context.Context, r *RenewAuthTokenRequest) (RenewAuthTokenResponse, error)
	RevokeAuthToken(ctx *context.Context, r *RevokeAuthTokenRequest) (RevokeAuthTokenResponse, error)
}

//GetAuthTokenRequest defines a request to the GetAuthToken Endpoint.
type GetAuthTokenRequest struct{}

//RenewAuthTokenRequest defines a request to the RenewAuthToken Endpoint.
type RenewAuthTokenRequest struct{}

//RevokeAuthTokenRequest defines a request to the RevokeAuthToken Endpoint.
type RevokeAuthTokenRequest struct{}

//GetAuthTokenResponse defines a response to the GetAuthToken Endpoint.
type GetAuthTokenResponse struct{}

//RenewAuthTokenResponse defines a response to the GetAuthToken Endpoint.
type RenewAuthTokenResponse struct{}

//RevokeAuthTokenResponse defines a response to RevokeAuthToken Endpoint
type RevokeAuthTokenResponse struct{}
