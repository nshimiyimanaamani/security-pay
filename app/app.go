package app

import (
	"context"
)

//App type is the root of the aoplication and implememts all the services
//defined by the different service interfaces. For that reason App is
//them amin container of App wide state.
type App struct{}

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

//RenewAuthToken receives an authentication request and returns a renewed jwt.
func (app *App) RenewAuthToken(ctx *context.Context, r *RenewTokenReq) (RenewTokenResp, error) {
	if err := r.validate(); err != nil {
		return RenewTokenResp{ID: r.ID}, err
	}
	return RenewTokenResp{
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
