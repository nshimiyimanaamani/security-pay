package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config contains application wide config
type Config struct {
	CloudRuntime string `validate:"required" envconfig:"PAYPACK_CLOUD_RUNTIME"`
	GoEnv        string `validate:"required" envconfig:"GO_ENV"`
	LogLevel     string `envconfig:"PAYPACK_LOG_LEVEL"` // http port
	Port         string `envconfig:"PAYPACK_HTTP_PORT"` // http port
	Secret       string `envconfig:"PAYPACK_SECRET"`
	DB           *DBConfig
	PaymentProxy *Payment
}

//default config
//	CloudRuntime: "none",
// GoEnv:        "development",
// LogLevel:     "debug",
// Port:         "8000",
// Secret:       "secret",
// DB:           &DBConfig{Endpoint: "postgres://postgres:password@localhost/test?sslmode=disable"},

// Load loads conf variables
func Load(prefix string) (*Config, error) {

	var c Config

	err := envconfig.Process(prefix, &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &c, nil
}
