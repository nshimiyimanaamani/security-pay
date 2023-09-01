package app

import (
	"context"
	"net/http"

	"github.com/nshimiyimanaamani/paypack-backend/backends/fdi"
	"github.com/nshimiyimanaamani/paypack-backend/backends/sms"
	"github.com/nshimiyimanaamani/paypack-backend/core/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/config"
	"github.com/quarksgroup/paypack-go/paypack/api"
	"github.com/quarksgroup/paypack-go/paypack/transport/oauth"
)

// InitPaymentClient initialises the payment gateway
func InitPaymentClient(ctx context.Context, cfg *config.Config) (payment.Client, error) {

	tr := &http.Client{
		Transport: &oauth.Transport{
			Scheme: oauth.SchemeBearer,
			Source: oauth.ContextTokenSource(),
			Base:   http.DefaultTransport,
		},
	}

	cli, err := api.New(cfg.Payment.PaymentURL, tr.Transport)
	if err != nil {
		return nil, err
	}

	return fdi.New(cli, cfg.Payment.AppID, cfg.Payment.Secret, cfg.GoEnv)
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
