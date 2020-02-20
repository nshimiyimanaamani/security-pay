package config

// SmsConfig ...
type SmsConfig struct {
	SmsURL   string `validate:"required" envconfig:"PAYPACK_SMS_URL"`
	SenderID string `validate:"required" envconfig:"PAYPACK_SMS_SENDER_ID"`
	Secret   string `validate:"required" envconfig:"PAYPACK_SMS_APP_SECRET"`
	AppID    string `validate:"required" envconfig:"PAYPACK_SMS_APP_ID"`
	//Callback string `validate:"required" envconfig:"PAYPACK_SMS_CALLBACK_URL"`
}
