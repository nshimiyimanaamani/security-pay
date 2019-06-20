package config

//Config defines the configuration getter API
type Config interface {
	FetchPrivateKey() string
	FetchPublicKey() string
}

var _Config = (*config)(nil)

//Config holds the application level configuration.
type config struct {
	PrivateKey string
	PublicKey  string
}

// New loads configuration
func New() (Config, error) {
	return &config{}, nil
}

//GetPrivateKey returns the private key.
func (cfg *config) FetchPrivateKey() string {
	return cfg.PrivateKey
}

//GetPublicKey returns the public key.
func (cfg *config) FetchPublicKey() string {
	return cfg.PublicKey
}
