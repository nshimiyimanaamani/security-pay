package app

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/backends/fdi"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
)

// InitPBackend initialises the payment gateway
func InitPBackend(ctx context.Context, cfg *config.PaymentConfig) (payment.Backend, error) {
	opts := &fdi.ClientOptions{
		URL:       cfg.PaymentURL,
		AppID:     cfg.AppID,
		AppSecret: cfg.Secret,
		Callback:  cfg.Callback,
	}
	return fdi.NewBackend(opts)
}
