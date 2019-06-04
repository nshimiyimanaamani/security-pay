package app

import (
	"testing"
	"context"
)

func TestGetAuthToken(t *testing.T){
	app:= New()

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

		///there should be an error if the request lacks an account fields
		{
			name:"request lacks 'account field'", 
			request: &GetTokenReq{ID:"1",Email:"example", Password:"pass"}, 
			token:false, 
			err: ErrorInvalidRequest,
		},

		{
			name:"request email is not valid", 
			request: &GetTokenReq{ID: "1", Email: "bad email", Password:"password", Account:"remera"}, 
			token:false, 
			err:ErrorInvalidEmail ,
		},

		{
			name:"request password is not valid", 
			request: &GetTokenReq{ID: "1", Email: "example", Password:"invalid password", Account:"remera"}, 
			token:false, 
			err: ErrorInvalidPassword,
		},

		{
			name:"request account does not exist", 
			request: &GetTokenReq{ID: "1", Email: "example", Password:"password", Account:"gasabo"}, 
			token:false, 
			err:ErrorAccountNotFound,
		},

		{
			name:"request has all the required", 
			request: &GetTokenReq{ID: "1", Email: "example", Password:"password", Account:"remera"}, 
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
	app:= New()

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
	app:= New()

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

func TestVerifyAuthToken(t *testing.T){
	app:= New()

	testcases:= []struct{
		name string
		request *VerifyTokenReq
		valid bool
		err error
	}{
		////there shoud be an error if the request doesn't the request id
		{
			name:"request lacks 'a request id'", 
			request: &VerifyTokenReq{Token:[]byte("token")}, 
			valid:false, 
			err: ErrorInvalidRequest,
		},

		{
			name:"request lacks a 'token field'", 
			request: &VerifyTokenReq{ID: "1"}, 
			valid:false, 
			err:ErrorInvalidRequest,
		},

		{
			name:"invalid token", 
			request: &VerifyTokenReq{ID: "1", Token:[]byte(" invalid token")}, 
			valid:false, 
			err:nil,
		},

		{
			name:"invalid token", 
			request: &VerifyTokenReq{ID: "1", Token:[]byte("valid token")}, 
			valid:true, 
			err:nil,
		},
	}

	for _,tc:=range testcases{
		t.Run(tc.name, func(t *testing.T){
			ctx:= context.Background()

			res, err:= app.VerifyAuthToken(&ctx, tc.request)

			if err!=tc.err{
				t.Errorf("expected error message to be '%v' got '%v'", tc.err, err)
			}

			//request and response must have the same id.
			if res.ID != tc.request.ID{
				t.Errorf("expected response and request to have the same id found-> req:%s | res:%s", tc.request.ID, res.ID)
			}

			//the response must have a token field if there is no error
			if tc.valid != res.Valid{
				t.Errorf("expected the valid field to be '%v' got '%v'", tc.valid, res.Valid)
			}
		})
	}
}