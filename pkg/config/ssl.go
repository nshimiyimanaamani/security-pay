package config

import validate "github.com/go-playground/validator/v10"

// SSLConfig contains the http server
type SSLConfig struct {
	Key         string
	Certificate string
}

func (conf *SSLConfig) Validate() error {
	validator := validate.New()
	return validator.Struct(conf)
}
