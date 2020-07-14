package fdi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rugwirobaker/paypack-backend/backends/encoding"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Timeout sets the default client timeout
const Timeout = 30 * time.Second

var _ payment.Client = (*client)(nil)

type client struct {
	URL       string
	ID        string
	Secret    string
	Token     string
	DCallback string
	CCallback string
	client    http.Client
}

// ClientOptions collects the NewBackend creattion options
type ClientOptions struct {
	URL       string
	AppSecret string
	AppID     string
	DCallback string
	CCallback string
}

// New ...
func New(opts *ClientOptions) payment.Client {
	if opts == nil {
		panic("fdi.Backend: absolutely unacceptable backend opts")
	}
	client := &client{
		URL:       opts.URL,
		DCallback: opts.DCallback,
		CCallback: opts.CCallback,
		ID:        opts.AppID,
		Secret:    opts.AppSecret,
		client:    http.Client{Timeout: Timeout},
	}
	return client
}

func (cli *client) Pull(ctx context.Context, tx payment.Payment) (payment.Response, error) {
	const op errors.Op = "backends/fdi/client.Pull"

	var empty payment.Response

	body := Request{
		TrxRef:      tx.ID,
		AccountID:   cli.ID,
		Msisdn:      tx.MSISDN,
		Amount:      tx.Amount,
		ChannelID:   string(tx.Method),
		CallbackURL: cli.DCallback,
	}

	raw, err := encoding.Serialize(body)
	if err != nil {
		return empty, errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, cli.URL+"/momo/pull", bytes.NewReader(raw))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// add authentication header
	token, err := cli.Auth()
	if err != nil {
		return empty, errors.E(op, err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := cli.client.Do(req)

	if err != nil {
		return empty, errors.E(op, err)
	}

	defer resp.Body.Close()

	res := &Response{}

	if err := encoding.Deserialize(resp.Body, res); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	status := payment.Response{
		Message: res.Message,
		Status:  res.Status,
		TxID:    res.Data.TrxRef,
		TxState: payment.State(res.Data.State),
	}
	return status, nil
}

func (cli *client) Push(ctx context.Context, tx payment.Payment) (payment.Response, error) {
	const op errors.Op = "backends/fdi/client.Push"

	var empty payment.Response

	body := Request{
		TrxRef:      tx.ID,
		AccountID:   cli.ID,
		Msisdn:      tx.MSISDN,
		Amount:      tx.Amount,
		ChannelID:   string(tx.Method),
		CallbackURL: cli.CCallback,
	}

	raw, err := encoding.Serialize(body)
	if err != nil {
		return empty, errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, cli.URL+"/momo/push", bytes.NewReader(raw))

	req.Header.Add("Accept", "application/json")

	// add authentication header
	token, err := cli.Auth()
	if err != nil {
		return empty, errors.E(op, err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := cli.client.Do(req)

	res := &Response{}

	if err := encoding.Deserialize(resp.Body, res); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	status := payment.Response{
		Message: res.Message,
		Status:  res.Status,
		TxID:    res.Data.TrxRef,
		TxState: payment.State(res.Data.State),
	}
	return status, nil
}

func (cli *client) Status(ctx context.Context) (int, error) {
	const op errors.Op = "backends/fdi/client.Status"

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
	if err := encoding.Deserialize(resp.Body, status); err != nil {
		return http.StatusInternalServerError, errors.E(op, err, errors.KindUnexpected)
	}

	return resp.StatusCode, nil
}

func (cli *client) Auth() (string, error) {
	const op errors.Op = "backends/fdi/client.Auth"

	// assemble request
	body := &authorization{ID: cli.ID, Secret: cli.Secret}

	bs, err := encoding.Serialize(body)
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

	if err := encoding.Deserialize(resp.Body, res); err != nil {
		return "", errors.E(op, err, errors.KindUnexpected)
	}

	if err := res.Validate(); err != nil {
		return "", errors.E(op, err, errors.KindUnexpected)
	}
	if token := res.Data.Token; token == "" {
		return "", errors.E(op, "unable to authenticate client", errors.KindUnexpected)
	}
	token := res.Data.Token

	return token, nil
}
