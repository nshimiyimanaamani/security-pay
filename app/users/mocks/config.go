package mocks

import (
	"github.com/rugwirobaker/paypack-backend/app/config"
)

var _ config.Config = (*configutationMock)(nil)

type configutationMock struct{}

//NewConfiguration creates "mirror" configuration provider
func NewConfiguration() config.Config {
	return &configutationMock{}
}

func (cfg *configutationMock) FetchPrivateKey() string {
	var key = `
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
	return key
}

func (cfg *configutationMock) FetchPublicKey() string {
	var key = `
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
	return key
}
