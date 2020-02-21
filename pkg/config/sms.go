package config

import validate "github.com/go-playground/validator/v10"

// SmsConfig ...
type SmsConfig struct {
	SmsURL   string `validate:"required" envconfig:"PAYPACK_SMS_APP_URL"`
	SenderID string `validate:"required" envconfig:"PAYPACK_SMS_SENDER_ID"`
	Secret   string `validate:"required" envconfig:"PAYPACK_SMS_APP_SECRET"`
	AppID    string `validate:"required" envconfig:"PAYPACK_SMS_APP_ID"`
	//Callback string `validate:"required" envconfig:"PAYPACK_SMS_CALLBACK_URL"`
}

func (conf *SmsConfig) Validate() error {
	validator := validate.New()
	return validator.Struct(conf)
}
