package app

import (
	"context"
)

//Application type is the root of the aoplication and implememts all the services
//defined by the different service interfaces. For that reason Application is
//them amin container of application wide state.
type Application struct{}

//GetAuthToken receives an authentication request and returns a jwt token.
func (app *Application) GetAuthToken(ctx *context.Context, r *GetTokenReq) (GetTokenResp, error) {
	return GetTokenResp{}, nil
}

//RenewAuthToken receives an authentication request and returns a renewed jwt.
func (app *Application) RenewAuthToken(ctx *context.Context, r *RenewTokenReq) (RenewTokenResp, error) {
	return RenewTokenResp{}, nil
}

//RevokeAuthToken receives an revokation request and revokes the given token.
func (app *Application) RevokeAuthToken(ctx *context.Context, r *RevokeTokenReq) (RevokeTokenResp, error) {
	return RevokeTokenResp{}, nil
}
