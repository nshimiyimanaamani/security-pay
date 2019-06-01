package app

import (
	"context"
	"errors"
)

//API wide error definitions
var (
	ErrorInvalidRequest = errors.New("invalid request")
)

//AuthService defines the authentication API.
type AuthService interface {
	GetAuthToken(ctx *context.Context, r *GetTokenReq) (GetTokenResp, error)
	RenewAuthToken(ctx *context.Context, r *RenewTokenReq) (RenewTokenResp, error)
	RevokeAuthToken(ctx *context.Context, r *RevokeTokenReq) (RevokeTokenResp, error)
}

//GetTokenReq defines a request to the GetAuthToken Endpoint.
type GetTokenReq struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//validate checks if the request is valid.
func (r *GetTokenReq) validate() error {
	if len(r.Email) < 1 || len(r.Password) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

//RenewTokenReq defines a request to the RenewAuthToken Endpoint.
type RenewTokenReq struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

func (r *RenewTokenReq) validate() error {
	if len(r.Token) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

//RevokeTokenReq defines a request to the RevokeAuthToken Endpoint.
type RevokeTokenReq struct {
	ID string `json:"id"`
}

//GetTokenResp defines a response to the GetAuthToken Endpoint.
type GetTokenResp struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

//RenewTokenResp defines a response to the GetAuthToken Endpoint.
type RenewTokenResp struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

//RevokeTokenResp defines a response to RevokeAuthToken Endpoint
type RevokeTokenResp struct {
	ID string `json:"id"`
}
