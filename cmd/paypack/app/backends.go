package app

import (
	"context"
	"net/http"

	"github.com/quarksgroup/paypack-go/paypack/api"
	"github.com/quarksgroup/paypack-go/paypack/transport/oauth"
	"github.com/rugwirobaker/paypack-backend/backends/fdi"
	"github.com/rugwirobaker/paypack-backend/backends/sms"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
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

	return fdi.New(cli, "8163b418-c3d7-11ed-9a56-dead64802bd2", "16423f218545d1cb49171764081c391eda39a3ee5e6b4b0d3255bfef95601890afd80709", cfg.GoEnv)
}

// InitSMSBackend ...
func InitSMSBackend(ctx context.Context, cfg *config.SMSConfig) (notifs.Backend, error) {
	opts := &sms.Options{
		URL:       "https://messaging.fdibiz.com/api/v1",
		SenderID:  "PayPack",
		AppID:     "228FCC98-027B-4BD4-84F5-D691E006B4E0",
		AppSecret: "4F7FE03E-0E34-44C6-9B41-8D8A6DDD2473",
	}
	return sms.New(opts)

}
