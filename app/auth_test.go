package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	app:= &Application{}

	testcases:= []struct{
		request 	*GetTokenReq
		err 		error
	}{
		{request: &GetTokenReq{ID: "1"}, err:nil},
	} 
	for _,tc:=range testcases{
		ctx:= context.Background()
		res, err:= app.GetAuthToken(&ctx, tc.request)

		if err!=tc.err{
			t.Errorf("expected error to be nil")
		}

		//request and response must have the same id.
		if res.ID != tc.request.ID{
			t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
		} 

		//the response must have a token field
		if len(res.Token)< 1{
			t.Errorf("the token field must not be empty")
		}
	}
}

func TestRenewAuthToken(t *testing.T){
	app:= &Application{}

	testcases:= []struct{
		request *RenewTokenReq
		err     error
	}{
		{request: &RenewTokenReq{ID: "1"}, err:nil},
	}

	for _,tc:= range testcases{
		ctx:= context.Background()

		res, err:= app.RenewAuthToken(&ctx, tc.request)

		if err!=tc.err{
			t.Errorf("expected error to be nil")
		}

		//request and response must have the same id.
		if res.ID != tc.request.ID{
			t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
		}

		//the response must have a token field
		if len(res.Token)< 1{
			t.Errorf("the token field must not be empty")
		}
	}
}

func TestRevokeAuthToken(t *testing.T){
	app:= &Application{}

	testcases:= []struct{
		request *RevokeTokenReq
		err     error
	}{
		{request: &RevokeTokenReq{ID: "1"}, err:nil},
	}

	for _,tc:= range testcases{
		ctx:= context.Background()

		res,err:= app.RevokeAuthToken(&ctx, tc.request)
		
		if err!=nil{
			t.Errorf("expected error to be nil")
		}

		//request and response must have the same id.
		if res.ID != tc.request.ID{
			t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
		} 
	}
}