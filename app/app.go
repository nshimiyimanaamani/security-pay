package app

import (
	"context"
)

//App type is the root of the aoplication and implememts all the services
//defined by the different service interfaces. For that reason App is
//them amin container of App wide state.
type App struct {
	Config *Config
}

//New initializes a new App instance.
func New() (*App, error) {
	app := &App{}
	cfg, err := NewConfig()
	if err != nil {
		return nil, err //improve code coverage
	}
	app.Config = cfg
	return app, nil
}

//GetAuthToken receives an authentication request and returns a jwt token.
//and otherwise a validation error and an empty response.
func (app *App) GetAuthToken(ctx *context.Context, r *GetTokenReq) (GetTokenResp, error) {
	if err := r.validate(); err != nil {
		return GetTokenResp{ID: r.ID}, err
	}
	if err := ValidateCredentials(r.Email, r.Password, r.Account); err != nil {
		return GetTokenResp{ID: r.ID}, err
	}
	return GetTokenResp{
		ID:    r.ID,
		Token: []byte("token"),
	}, nil
}

//RefreshAuthToken receives an authentication request and returns a Refreshed jwt.
//and otherwise a validation error and an empty response.
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
//and otherwise a validation error and an empty response.
func (app *App) RevokeAuthToken(ctx *context.Context, r *RevokeTokenReq) (RevokeTokenResp, error) {
	if err := r.validate(); err != nil {
		return RevokeTokenResp{ID: r.ID}, err
	}
	return RevokeTokenResp{ID: r.ID}, nil
}

//VerifyAuthToken receives a a verification request anf returns true if the given token is valid
//and otherwise a validation error and an empty response.
func (app *App) VerifyAuthToken(ctx *context.Context, r *VerifyTokenReq) (VerifyTokenResp, error) {
	if err := r.validate(); err != nil {
		return VerifyTokenResp{ID: r.ID}, err
	}
	return VerifyTokenResp{ID: r.ID, Valid: ValidateToken(r.Token)}, nil
}

/**
 * @todo Add JWT TOKEN
 * @body I should be adding a real Jwt token by commiting the code in models ounce done
 */
