package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &GetTokenReq{
		ID: "1",
	}

	res, err:= app.GetAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}

	//request and response must have the same id.
	if res.ID != req.ID{
		t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", req.ID, res.ID)
	} 
}

func TestRenewAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &RenewTokenReq{
		ID: "10",
	}

	res, err:= app.RenewAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
	
	//request and response must have the same id.
	if res.ID != req.ID{
		t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", req.ID, res.ID)
	} 
}

func TestRevokeAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}
	req:= &RevokeTokenReq{
		ID: "10",
	}

	res,err:= app.RevokeAuthToken(&ctx, req)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}

	//request and response must have the same id.
	if res.ID != req.ID{
		t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", req.ID, res.ID)
	} 
}