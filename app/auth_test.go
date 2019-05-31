package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &GetAuthTokenRequest{}

	err:= app.GetAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
}

func TestRenewAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &RenewAuthTokenRequest{}

	err:= app.RenewAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
}

func TestRevokeAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &RevokeAuthTokenRequest{}

	err:= app.RevokeAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
}