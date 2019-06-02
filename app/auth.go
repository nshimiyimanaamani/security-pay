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
	RefreshAuthToken(ctx *context.Context, r *RefreshTokenReq) (RefreshTokenResp, error)
	RevokeAuthToken(ctx *context.Context, r *RevokeTokenReq) (RevokeTokenResp, error)
}

//GetTokenReq defines a request to the GetAuthToken Endpoint.
type GetTokenReq struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Account  string `json:"account"`
}

//validate checks if the request is valid.
func (r *GetTokenReq) validate() error {
	if len(r.ID) < 1 || len(r.Email) < 1 || len(r.Password) < 1 || len(r.Account) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

//RefreshTokenReq defines a request to the RefreshAuthToken Endpoint.
type RefreshTokenReq struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

func (r *RefreshTokenReq) validate() error {
	if len(r.ID) < 1 || len(r.Token) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

//RevokeTokenReq defines a request to the RevokeAuthToken Endpoint.
type RevokeTokenReq struct {
	ID string `json:"id"`
}

func (r *RevokeTokenReq) validate() error {
	if len(r.ID) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

//GetTokenResp defines a response to the GetAuthToken Endpoint.
type GetTokenResp struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

//RefreshTokenResp defines a response to the GetAuthToken Endpoint.
type RefreshTokenResp struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

//RevokeTokenResp defines a response to RevokeAuthToken Endpoint
type RevokeTokenResp struct {
	ID string `json:"id"`
}
