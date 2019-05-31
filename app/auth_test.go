package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &GetTokenReq{}

	res, err:= app.GetAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
	if len(res.ID) == 0{
		t.Errorf("response id must be different from nil")
	} 
}

func TestRenewAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &RenewTokenReq{}

	res, err:= app.RenewAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
	if len(res.ID) == 0{
		t.Errorf("response id must be different from nil")
	} 
}

func TestRevokeAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &RevokeTokenReq{}

	res,err:= app.RevokeAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
	if len(res.ID) == 0{
		t.Errorf("response id must be different from nil")
	} 
}