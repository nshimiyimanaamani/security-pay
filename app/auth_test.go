package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	app:= NewApp()

	testcases:= []struct{
		name        string
		request 	*GetTokenReq
		token       bool
		err 		error
		
	}{	

		////there shoud be an error if the request doesn't the request id
		{
			name:"request lacks 'a request id'", 
			request: &GetTokenReq{Email: "example", Password:"pass"}, 
			token:false, err: ErrorInvalidRequest,
		},
		////there shoud be an error if the request doesn't have any of the required fields
		{
			name:"request lacks 'all the fields'", 
			request: &GetTokenReq{ID:"1"}, 
			token:false, 
			err: ErrorInvalidRequest,
		},

		///there should be an error if the request lacks email field
		{
			name:"request lacks 'email field'", 
			request: &GetTokenReq{ID:"1",Password:"pass"}, 
			token:false,
			err: ErrorInvalidRequest,
		},

		///there should be an error if the request lacks a password fields
		{
			name:"request lacks 'password field'", 
			request: &GetTokenReq{ID:"1",Email:"example"}, 
			token:false, 
			err: ErrorInvalidRequest,
		},

		///perfect request
		{
			name:"request has all the required", 
			request: &GetTokenReq{ID: "1", Email: "example", Password:"pass"}, 
			token:true, 
			err:nil,
		},
	}
	for _,tc:=range testcases{
		t.Run(tc.name, func(t *testing.T){
			ctx:= context.Background()
			res, err:= app.GetAuthToken(&ctx, tc.request)

			//request and response must have the same id.
			if res.ID != tc.request.ID{
				t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
			} 

			if err!=tc.err{
				t.Errorf("expected error message to be '%v' got '%v'", tc.err, err)
			}

			//the response must have a token field if there is no error
			if tc.token == true && len(res.Token)< 1{
				t.Errorf("the token field must not be empty")
			}
		})
	}
}

func TestRovken(t *testing.T){
	app:= NewApp()

	testcases:= []struct{
		name    string
		request *RefreshTokenReq
		token   bool
		err     error
	}{	
		////there shoud be an error if the request doesn't the request id
		{
			name:"request lacks 'a request id'", 
			request: &RefreshTokenReq{Token:[]byte("old token")}, 
			token:false, 
			err: ErrorInvalidRequest,
		},

		{
			name:"request lacks a 'token field'", 
			request: &RefreshTokenReq{ID: "1"}, 
			token:false, 
			err:ErrorInvalidRequest,
		},

		{
			name:"request has 'all the fields'", 
			request: &RefreshTokenReq{ID: "1", 
			Token:[]byte("old token")}, 
			token:true, 
			err:nil,
		},
	}

	for _,tc:= range testcases{
		t.Run(tc.name, func(t *testing.T){
			ctx:= context.Background()

			res, err:= app.RefreshAuthToken(&ctx, tc.request)

			if err!=tc.err{
				t.Errorf("expected error message to be '%v' got '%v'", tc.err, err)
			}

			//request and response must have the same id.
			if res.ID != tc.request.ID{
				t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
			}

			//the response must have a token field if there is no error
			if tc.token == true && len(res.Token)< 1{
				t.Errorf("the token field must not be empty")
			}
		})
	}
}

func TestRevokeAuthToken(t *testing.T){
	app:= NewApp()

	testcases:= []struct{
		name 	string
		request *RevokeTokenReq
		err     error
	}{
		{name:"valid request", request: &RevokeTokenReq{ID: "1"}, err:nil},

		{name:"request lacks 'id'", request: &RevokeTokenReq{}, err:ErrorInvalidRequest},
	}

	for _,tc:= range testcases{
		t.Run(tc.name, func(t *testing.T){
			ctx:= context.Background()

			res,err:= app.RevokeAuthToken(&ctx, tc.request)
		
			if err!=tc.err{
				t.Errorf("expected error message to be '%v' got '%v'", tc.err, err)
			}

			//request and response must have the same id.
			if res.ID != tc.request.ID{
				t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
			} 
		})
	}
}