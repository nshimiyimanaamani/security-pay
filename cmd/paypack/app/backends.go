package app

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/backends/fdi"
	"github.com/rugwirobaker/paypack-backend/backends/sms"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
)

// InitPaymentClient initialises the payment gateway
func InitPaymentClient(ctx context.Context, cfg *config.PaymentConfig) payment.Client {
	opts := &fdi.ClientOptions{
		URL:       cfg.PaymentURL,
		AppID:     cfg.AppID,
		AppSecret: cfg.Secret,
		DCallback: cfg.DCallback,
		CCallback: cfg.CCallback,
	}
	return fdi.New(opts)
}

// InitSMSBackend ...
func InitSMSBackend(ctx context.Context, cfg *config.SMSConfig) (notifs.Backend, error) {
	opts := &sms.Options{
		URL:       cfg.SmsURL,
		SenderID:  cfg.SenderID,
		AppID:     cfg.AppID,
		AppSecret: cfg.Secret,
	}
	return sms.New(opts)

}
