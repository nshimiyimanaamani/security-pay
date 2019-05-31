package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}

	err:= app.GetAuthToken(&ctx)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
}

func TestRenewAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}

	err:= app.RenewAuthToken(&ctx)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
}

func TestRevokeAuthToken(t *testing.T){
	ctx:= context.Background()

	app:= &Application{}

	err:= app.RevokeAuthToken(&ctx)
	if err!=nil{
		t.Errorf("expected error to be nil")
	}
}