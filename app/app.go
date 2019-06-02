package app

import (
	"context"
)

//App type is the root of the aoplication and implememts all the services
//defined by the different service interfaces. For that reason App is
//them amin container of App wide state.
type App struct{}

//New initializes a new App instance.
func New() *App {
	return &App{}
}

//GetAuthToken receives an authentication request and returns a jwt token.
func (app *App) GetAuthToken(ctx *context.Context, r *GetTokenReq) (GetTokenResp, error) {
	if err := r.validate(); err != nil {
		return GetTokenResp{ID: r.ID}, err
	}
	return GetTokenResp{
		ID:    r.ID,
		Token: []byte("token"),
	}, nil
}

//RefreshAuthToken receives an authentication request and returns a Refreshed jwt.
func (app *App) RefreshAuthToken(ctx *context.Context, r *RefreshTokenReq) (RefreshTokenResp, error) {
	if err := r.validate(); err != nil {
		return RefreshTokenResp{ID: r.ID}, err
	}
	return RefreshTokenResp{
		ID:    r.ID,
		Token: []byte("new token"),
	}, nil
}

//RevokeAuthToken receives an revokation request and revokes the given token.
func (app *App) RevokeAuthToken(ctx *context.Context, r *RevokeTokenReq) (RevokeTokenResp, error) {
	if err := r.validate(); err != nil {
		return RevokeTokenResp{ID: r.ID}, err
	}
	return RevokeTokenResp{ID: r.ID}, nil
}

//VerifyAuthToken receives a a verification request anf returns true if the given token is valid
func (app *App) VerifyAuthToken(ctx *context.Context, r *VerifyTokenReq) (VerifyTokenResp, error) {
	if err := r.validate(); err != nil {
		return VerifyTokenResp{ID: r.ID}, err
	}
	return VerifyTokenResp{ID: r.ID, Valid: true}, nil
}
