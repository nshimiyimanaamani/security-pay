package app

import (
	"crypto/rand"
	"crypto/rsa"
)

//GetConfigKey defines the configuration getter API
type GetConfigKey interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
}

//Config holds the application level configuration.
type Config struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewConfig loads configuration
func NewConfig() (*Config, error) {
	bitSize := 2048

	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	public := &key.PublicKey

	return &Config{
		PrivateKey: key,
		PublicKey:  public,
	}, nil
}

//GetPrivateKey returns the private key.
func (cfg *Config) GetPrivateKey() *rsa.PrivateKey {
	return cfg.PrivateKey
}

//GetPublicKey returns the public key.
func (cfg *Config) GetPublicKey() *rsa.PublicKey {
	return cfg.PublicKey
}
