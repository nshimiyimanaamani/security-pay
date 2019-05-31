package app

import (
	"context"
)

//AuthService defines the authentication API.
type AuthService interface {
	GetAuthToken(ctx *context.Context, r *GetTokenReq) (GetTokenResp, error)
	RenewAuthToken(ctx *context.Context, r *RenewTokenReq) (RenewTokenResp, error)
	RevokeAuthToken(ctx *context.Context, r *RevokeTokenReq) (RevokeTokenResp, error)
}

//GetTokenReq defines a request to the GetAuthToken Endpoint.
type GetTokenReq struct{}

//RenewTokenReq defines a request to the RenewAuthToken Endpoint.
type RenewTokenReq struct{}

//RevokeTokenReq defines a request to the RevokeAuthToken Endpoint.
type RevokeTokenReq struct{}

//GetTokenResp defines a response to the GetAuthToken Endpoint.
type GetTokenResp struct {
	ID string `json:"id"`
}

//RenewTokenResp defines a response to the GetAuthToken Endpoint.
type RenewTokenResp struct {
	ID string `json:"id"`
}

//RevokeTokenResp defines a response to RevokeAuthToken Endpoint
type RevokeTokenResp struct {
	ID string `json:"id"`
}
