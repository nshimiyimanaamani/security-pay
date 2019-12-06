package app

import (
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/payment/fdi"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
)

// CreatePaymentBackend initialises the payment gateway
func CreatePaymentBackend(cfg *config.PaymentConfig) (payment.Backend, error) {
	opts := &fdi.ClientOptions{
		URL:       cfg.PaymentURL,
		AppID:     cfg.AppID,
		AppSecret: cfg.Secret,
		Callback:  cfg.Callback,
	}
	b := fdi.NewBackend(opts)
	return b, nil
}
