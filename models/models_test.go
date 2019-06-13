package models

import (
	"crypto/rsa"
	"testing"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestGenerate(t *testing.T){
	//token:= &Token{}

	privKey, err:= jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err!=nil{
		t.Fatalf("unexpected error '%v'",err)
	}

	claims:= &Claims{
		Account: "Remera",
		Email: "example@!gmail.com",
		Admin: true,
	}

	token, err:= Generate(privKey, claims)

	if err!=nil && len(token)<1{
		t.Fatalf("unexpected error '%v'",err)
	}

	if !validate(token, claims, privKey, t){
		t.Errorf("token string should be valid")
	}
}

func TestParse(t *testing.T){
	key, err:= jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err!=nil{
		t.Fatalf("unexpected error '%v'",err)
	}

	tokenString:= "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50IjoiUmVtZXJhIiwiRW1haWwiOiJleGFtcGxlQCFnbWFpbC5jb20iLCJBZG1pbiI6dHJ1ZX0.lzIaopc4xBIR276REB8nIVqnszULc2dK7l_pdIOe8_rmOU9cjl4Ibp_kgwczXpGyIXYVuJq-uxO4C8EjED5n5XKkKTP19oa6t-fJyR0aXu_6CRJaAdp8btYhyxIXuWi8kq8MTEZZFr03CS7Wbaf7IRF9bTScaISTm5T-doOjtivUWsIrO7dlai7ddPG2YApV5mv0IgaESn51YV6cSglWVk7xfDoV7pPHElf8kiLUz3VdKWriYjHoFuWc_JmU9cWPuwmPwl4JoUOcxXGWi6uJz3T_yL78-zDYDXQkWEmba5yoxVUnWj22HiRo3NUlyI4uYBSZCb70Q78HgCX5DbaziQ"

	token, err:= Parse(tokenString, key)
	if err!=nil{
		t.Fatalf("unexpected error '%v'",err)
	}

	if !token.Validate(){
		t.Errorf("expected token to valid")
	}


 }

func validate(tokenString string, claims *Claims, priv *rsa.PrivateKey, t *testing.T)bool{
	t.Helper()

	pub:= priv.Public()
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return pub, nil
	})
	
	if err!=nil{
		t.Fatalf("unexpected error '%v'", err)
	}
	return token.Valid
}

var privateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA5k7M61N3NYkkDMuNhfd2CZnd5Gxu6rpzBgYkbFKEP2koGZv4
a6arQyW/eSL5bfFdzOdPcq3NKOpWxOgoMgLkI1QgI56orMGjAhKzUiN8NPtNFQcf
S19R+OIvikrXfO8Z7tovRIyFbKYx/Ar5v7Q5ZH6mhOcPmSH7tzwqb7T9VCGDiml/
JwKW+ErcqZl0PkHO6855BMA4gyiHHRwnCwYScrH45glnTyxTfAFizwqfOwzMzQFN
LvsH3L3UfR7YQaEW/buISVdwA39CIoRYHza/Z9qpzmGXowICHCOzwyOUmFQR57cu
STC/ajcu3pVADv1Ur2sjxXyehdYBGeKFx3hS5wIDAQABAoIBABZX0r2JzWjeMyci
oEo85bCswsAkXOZczEfrVKFFqBrWwtMpNIKNGtRa1yaTZAtsfSMh1a1UezDa+ywD
MdMYQLXEtZF/FPIdnwjWc5smYihpsOK3XCvdxYAVwXLzK9CtCaEIfclysIcH4JWJ
Iw2cGG1NdC40lGjQyTDPn3ZS4rjEkGbUzpKrjqsnoh3qRUGu/VyU/9eWN5TrztbR
J9g7dEncU/VSWLIx1IUEd+Ay5sgEytltS4nGmteuQYhuyup7KeQpsOvZFUoQ7nZ3
NnfcwaEs00sgqvbk5FgxFpvGczPqNM8I1+r+ayxZ+Qmx2vkS8iADf7Ze2ZLOWsQ8
8akEfCECgYEA/9nOzEbcyPMUpUZrRzEmJrBUSTdveuOPqobZRUxXVpMiJIqpoguq
0JALZL6G3B8MbB8gkSKa2ZLv0fEskWaPOwIVzmfL7D9MUMCNQVHB3nv8c5xQ4fVA
HlhI7ETvc4PL2UMIy6WLjsy+O3PIShWe7R/TpXQYcdCQRduRkYAIrCkCgYEA5nEu
AmDTsT42upYjqPyCXpDnSSUWyy/ZIg7QBYYuB0aIbZ3g22Mh5TZ9kQCJ7GIor/2s
/xXfT3RcLKKkbpXkgbzaG0It0Ha4wooU2NU38GK1r+A/+edr6nxeQAJjhWdrcYbm
UHJEjIwxd1YRn+6paKZQ4YjFsmYyQr+8tJUj6I8CgYEAtdokqQG9MH/GrursmX+P
tHQklJ34eQqCNR0AFcd7VKfj3sFIbUuJsBCSaJsb1B2lgLxnM4G7Oua72ydnHDof
mDuVME6KnXMoVUVnoYPxHqhV+f6jZtghKPBrdLRS1nJZVCXXfJhAJ9HTbQKQ3Eed
3MGAd9ua/FrYES9NunOctnECgYB1ZG7V01HEVzc2MkoUSh534kWQo45LECMDEJy0
U7ibCDlz7hugZ43a4Lly5t1cSF0F2qsIf7H1HgfezTQLCd0Qoo5RmJMSQYi5wfIA
zA3lLcP0xr6Qpm35VEYHQbBFQ3wep2Qo0y1MlBaW/oeX+9Lddux5GF3uFdXA30BY
lilmzQKBgEs37af1oGYQVgcqGdqvggIxKSMTRQePd/jcvmApBd6PxbID/SNnupaN
B6RfoVnahqgAn/3v0mGVBPRkrduAyb+I+ZaFz6sbrvj+Uv6qq//k5IX3jsgG9bYv
kr/6HUjbnnxpRv7ojjnTGPn470EaiqcLkrcg6CS7/N+EUKaoW9Cu
-----END RSA PRIVATE KEY-----
`

var publicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5k7M61N3NYkkDMuNhfd2
CZnd5Gxu6rpzBgYkbFKEP2koGZv4a6arQyW/eSL5bfFdzOdPcq3NKOpWxOgoMgLk
I1QgI56orMGjAhKzUiN8NPtNFQcfS19R+OIvikrXfO8Z7tovRIyFbKYx/Ar5v7Q5
ZH6mhOcPmSH7tzwqb7T9VCGDiml/JwKW+ErcqZl0PkHO6855BMA4gyiHHRwnCwYS
crH45glnTyxTfAFizwqfOwzMzQFNLvsH3L3UfR7YQaEW/buISVdwA39CIoRYHza/
Z9qpzmGXowICHCOzwyOUmFQR57cuSTC/ajcu3pVADv1Ur2sjxXyehdYBGeKFx3hS
5wIDAQAB
-----END PUBLIC KEY-----
`