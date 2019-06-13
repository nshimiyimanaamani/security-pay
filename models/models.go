package models

import "crypto/rsa"

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//Token defines the jwt token
type Token struct {
	Token *jwt.Token
}

//Validate returns true if the token is valid or false otherwise
func (tk *Token) Validate() bool {
	return tk.Token.Valid
}

//Claims a struct that will be encoded to a JWT, embedds the jwt type.
type Claims struct {
	Account string
	Email   string
	Admin   bool
	jwt.StandardClaims
}

//Generate takes an rsa Private Key and Claims to return a new token and an non error.
func Generate(key *rsa.PrivateKey, claims *Claims) (string, error) {
	token := jwt.NewWithClaims((jwt.GetSigningMethod("RS256")), claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//Parse converts a string into  type Token.
//Takes a token String an returns a Token and a nil error ig the token string is valid.
/**
 * @todo Add a token Parseer function
 * @body The payment History interfaces deals with listing all the payment receipts.
 */
func Parse(tokenString string, key *rsa.PublicKey) (*Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return &Token{Token: &jwt.Token{}}, err
	}

	return &Token{
		Token: token,
	}, nil
}
