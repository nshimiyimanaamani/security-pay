package models

import "crypto/rsa"

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//Token defines the jwt token
type Token struct {
	Token string
}

//Claims a struct that will be encoded to a JWT, embedds the jwt type.
type Claims struct {
	Account string
	Email   string
	Admin   bool
	jwt.StandardClaims
}

//Generate takes an rsa Private Key and Claims to return a new token and an non error.
func Generate(key *rsa.PrivateKey, claims *Claims) (*Token, error) {
	token := jwt.NewWithClaims((jwt.GetSigningMethod("RS256")), claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &Token{Token: tokenString}, nil
}

//Parse converts a string into  type Token.
//Takes a token String an returns a Token and a nil error ig the token string is valid.
/**
 * @todo Add a token Parseer function
 * @body The payment History interfaces deals with listing all the payment receipts.
 */
func Parse(tkString string, key *rsa.PrivateKey) (*Token, error) {
	return nil, nil
}

//Validate verifies the validity a jwt token
//and retruns a non nil error if any criterion is not met
/**
 * @todo Add Payment History Interface
 * @body The payment History interfaces deals with listing all the payment receipts.
 */
func (tk *Token) Validate() error { return nil }
