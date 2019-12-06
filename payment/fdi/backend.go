package fdi

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Timeout sets the default client timeout
const Timeout = 30 * time.Second

type backend struct {
	URL       string
	AppSecret string
	AppID     string
	Token     string
	Callback  string
	client    http.Client
}

// ClientOptions collects the NewBackend creattion options
type ClientOptions struct {
	URL       string
	AppSecret string
	AppID     string
	Callback  string
}

// NewBackend ...
func NewBackend(opts *ClientOptions) payment.Backend {
	if opts == nil {
		panic("fdi.Backend: absolutely unacceptable backend opts")
	}
	client := http.Client{Timeout: Timeout}

	return &backend{
		URL:       opts.URL,
		AppSecret: opts.AppSecret,
		Callback:  opts.Callback,
		client:    client,
	}
}

func (cli *backend) Pull(ctx context.Context, tx payment.Transaction) (string, error) {
	const op errors.Op = "fdi.Pull"

	body := pullRequest{
		TrxRef:      tx.ID,
		AccountID:   cli.AppID,
		Msisdn:      tx.Phone,
		Amount:      tx.Amount,
		ChannelID:   tx.Method,
		CallbackURL: cli.Callback,
	}

	bits, err := Serialize(body)
	if err != nil {
		return "", errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, cli.URL+"/momo/pull", bytes.NewReader(bits))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	token := "Bearer" + cli.Token
	req.Header.Add("Authorization", "Bearer"+token)

	resp, err := cli.client.Do(req)

	if err != nil {
		return "", errors.E(op, err)
	}

	defer resp.Body.Close()

	var res pullResponse

	if err := Deserialize(resp.Body, &res); err != nil {
		return "", errors.E(op, err, errors.KindUnexpected)
	}

	if res.Status == "" {
		return "", errors.E(op, err, errors.KindUnexpected)
	}
	return res.Status, nil
}

func (cli *backend) Status(ctx context.Context) (int, error) {
	const op errors.Op = "fdi.Status"

	req, err := http.NewRequest(http.MethodGet, cli.URL+"/status", nil)
	if err != nil {
		return http.StatusInternalServerError, errors.E(op, err, errors.KindUnexpected)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := cli.client.Do(req)

	if err != nil {
		return http.StatusInternalServerError, errors.E(op, err, errors.KindUnexpected)
	}

	defer resp.Body.Close()

	var status statusResponse
	if err := Deserialize(resp.Body, &status); err != nil {
		return http.StatusInternalServerError, errors.E(op, err, errors.KindUnexpected)
	}

	return resp.StatusCode, nil
}

func (cli *backend) Auth(ctx context.Context) error {
	const op errors.Op = "fdi.Auth"

	body := authRequest{AppID: cli.AppID, Secret: cli.AppSecret}

	bs, err := Serialize(body)
	if err != nil {
		return errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, cli.URL+"/auth", bytes.NewReader(bs))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := cli.client.Do(req)

	if err != nil {
		return errors.E(op, err)
	}

	defer resp.Body.Close()

	var auth authResponse
	if err := Deserialize(resp.Body, &auth); err != nil {
		return errors.E(op, err)
	}

	token := auth.Data.Token
	if token == "" {
		return errors.E(op, "authentication failed", errors.KindNotFound)
	}

	return nil
}
