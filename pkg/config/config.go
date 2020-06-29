package config

import (
	"log"

	validate "github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Config contains application wide config
type Config struct {
	CloudRuntime string `validate:"required" envconfig:"PAYPACK_CLOUD_RUNTIME"`
	GoEnv        string `validate:"required" envconfig:"GO_ENV"`
	LogLevel     string `envconfig:"PAYPACK_LOG_LEVEL"` // http port
	Port         string `envconfig:"PORT"`              // http port
	Secret       string `envconfig:"PAYPACK_SECRET"`
	Postgres     *PostgresConfig
	Redis        *RedisConfig
	Payment      *PaymentConfig
	SMS          *SMSConfig
	USSD         *USSDConfig
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

// Validate configuration
func Validate(conf *Config) error {
	const op errors.Op = "pkg/config/Config.Validate"

	var validator = validate.New()

	if err := validator.Struct(conf); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
