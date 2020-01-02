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
	URL      string
	AppID    string
	Token    string
	Callback string
	client   http.Client
}

// ClientOptions collects the NewBackend creattion options
type ClientOptions struct {
	URL       string
	AppSecret string
	AppID     string
	Callback  string
}

// NewBackend ...
func NewBackend(opts *ClientOptions) (payment.Backend, error) {
	const op errors.Op = "backend.fdi.NewBackend"
	if opts == nil {
		panic("fdi.Backend: absolutely unacceptable backend opts")
	}
	client := http.Client{Timeout: Timeout}

	backend := &backend{
		URL:      opts.URL,
		Callback: opts.Callback,
		client:   client,
	}
	token, err := backend.Auth(opts.AppID, opts.AppSecret)
	if err != nil {
		return nil, errors.E(op, err)
	}
	backend.Token = token
	return backend, nil
}

func (cli *backend) Pull(ctx context.Context, tx payment.Transaction) (payment.Status, error) {
	const op errors.Op = "backend.fdi.Pull"

	empty := payment.Status{}

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
		return empty, errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, cli.URL+"/momo/pull", bytes.NewReader(bits))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	token := "Bearer " + cli.Token
	req.Header.Add("Authorization", token)

	resp, err := cli.client.Do(req)

	if err != nil {
		return empty, errors.E(op, err)
	}

	defer resp.Body.Close()

	res := &pullResponse{}

	if err := Deserialize(resp.Body, res); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	status := payment.Status{
		Message: res.Message,
		Status:  res.Status,
		TxID:    res.Data.TrxRef,
		TxState: payment.State(res.Data.State),
	}
	return status, nil
}

func (cli *backend) Status(ctx context.Context) (int, error) {
	const op errors.Op = "backend.fdi.Status"

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

	var status = &statusResponse{}
	if err := Deserialize(resp.Body, status); err != nil {
		return http.StatusInternalServerError, errors.E(op, err, errors.KindUnexpected)
	}

	return resp.StatusCode, nil
}

func (cli *backend) Auth(appID, appSecret string) (string, error) {
	const op errors.Op = "backends.fdi.Auth"

	// assemble request
	body := &authRequest{AppID: appID, Secret: appSecret}

	bs, err := Serialize(body)
	if err != nil {
		return "", errors.E(op, err)
	}
	req, err := http.NewRequest(http.MethodPost, cli.URL+"/auth", bytes.NewReader(bs))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// make request
	resp, err := cli.client.Do(req)
	if err != nil {
		return "", errors.E(op, err, errors.KindUnexpected)
	}

	// read response
	res := &authResponse{}

	if err := Deserialize(resp.Body, res); err != nil {
		return "", errors.E(op, err, errors.KindUnexpected)
	}
	if token := res.Data.Token; token == "" {
		return "", errors.E(op, "unable to authenticate client", errors.KindUnexpected)
	}
	return res.Data.Token, nil
}
