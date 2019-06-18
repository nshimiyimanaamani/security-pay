package app

/**
 * @todo Refactor the AuthService
 * @body The refactor service must retrun a model object [Token] which is part of the business logic.
 */

import (
	"bytes"
	"context"
	"errors"
)

//API wide error definitions
var (
	ErrorInvalidRequest  = errors.New("invalid request")
	ErrorInvalidEmail    = errors.New("email is not valid")
	ErrorInvalidPassword = errors.New("password is not valid")
	ErrorAccountNotFound = errors.New("the account does not exists")
)

//AuthService defines the authentication API.
type AuthService interface {
	//GetToken creates a new token.
	//It takes a context and 2 strings[email, account] and a boolean.
	//A succeful operation retuns a valid jwt token and nil error.
	GetToken(*context.Context, string, string, bool) (string, error)

	//RefreshToken renews an expired token.
	//It takes a context object and a string[token].
	//A succeful operation retuns a valid jwt token and nil error.
	RefreshToken(*context.Context, string) (string, error)

	//Revoke token invalidates a token.
	//It takes a context and a string[token]
	//A succeful operation retuns a nil error.
	RevokeToken(*context.Context, string) error

	//VerifyToken chechks wether a token is valid.
	//It takes a context and a string[token]
	//A succeful operation retuns a nil error.
	VerifyToken(*context.Context, string) error
}

//GetTokenReq defines a request to the GetAuthToken Endpoint.
type GetTokenReq struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Account  string `json:"account"`
}

//RefreshTokenReq defines a request to the RefreshAuthToken Endpoint.
type RefreshTokenReq struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
}

//RevokeTokenReq defines a request to the RevokeAuthToken Endpoint.
type RevokeTokenReq struct {
	ID string `json:"id"`
}

//VerifyTokenReq defines a request to the VerifyAuthToken Endpoint
type VerifyTokenReq struct {
	ID    string `json:"id"`
	Token []byte `json:"token"`
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

//VerifyTokenResp defines a response from the VerifyAuthToken Endpoint.
type VerifyTokenResp struct {
	ID    string `json:"id"`
	Valid bool   `json:"valid"`
}

//validate checks if the request is valid.
func (r *GetTokenReq) validate() error {
	if len(r.ID) < 1 || len(r.Email) < 1 || len(r.Password) < 1 || len(r.Account) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

func (r *RefreshTokenReq) validate() error {
	if len(r.ID) < 1 || len(r.Token) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

func (r *RevokeTokenReq) validate() error {
	if len(r.ID) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

func (r *VerifyTokenReq) validate() error {
	if len(r.ID) < 1 || len(r.Token) < 1 {
		return ErrorInvalidRequest
	}
	return nil
}

//ValidateToken validates a given token against a set of preset rules.
//Tt returns true if token is valid and otherwise false
func ValidateToken(token []byte) bool {
	return bytes.Equal(token, []byte("valid token"))
}

//ValidateCredentials returns true if all the given credentials are valid.
func ValidateCredentials(email, password, account string) error {
	validEmail := "example"
	validPassword := "password"
	validAccount := "remera"

	if email != validEmail {
		return ErrorInvalidEmail
	}
	if password != validPassword {
		return ErrorInvalidPassword
	}
	if account != validAccount {
		return ErrorAccountNotFound
	}
	return nil
}
