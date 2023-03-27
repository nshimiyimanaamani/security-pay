package sms

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/rugwirobaker/paypack-backend/backends/encoding"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Timeout sets the default client timeout
const Timeout = 30 * time.Second

var _ (notifs.Backend) = (*Backend)(nil)

// Backend ..
type Backend struct {
	URL      string
	Token    string
	Refresh  string
	SenderID string
	client   *http.Client
}

// Options ...
type Options struct {
	URL       string
	AppID     string
	AppSecret string
	SenderID  string
}

// New ...
func New(opts *Options) (*Backend, error) {
	const op errors.Op = "backends/sms/New"

	if opts == nil {
		panic("fdi.Backend: absolutely unacceptable backend opts")
	}

	client := &http.Client{Timeout: Timeout}

	backend := &Backend{
		URL:      opts.URL,
		client:   client,
		SenderID: opts.SenderID,
	}
	if err := backend.authorize(opts.AppID, opts.AppSecret); err != nil {
		return nil, errors.E(op, err)
	}
	return backend, nil
}

// Send ...
func (sms *Backend) Send(ctx context.Context, id, message string, recipients []string) error {
	const op errors.Op = "backend/sms/backend.Send"

	switch len(recipients) {
	case 0:
		return errors.E(op, "cannot send message to recipients", errors.KindBadRequest)
	case 1:
		err := sms.sendSingle(ctx, id, recipients[0], message)
		if err != nil {
			return errors.E(op, err)
		}
	default:
		err := sms.sendBulk(ctx, id, recipients, message)

		if err != nil {
			return errors.E(op, err)
		}
	}
	return nil
}

func (sms *Backend) sendSingle(ctx context.Context, id, recipient, message string) error {
	const op errors.Op = "backends/sms/Backend.sendSingle"

	reqBody := &singleMSISDNRequest{
		MSISDN:    recipient,
		Message:   message,
		SenderID:  "PayPack",
		Reference: id,
	}

	b, err := encoding.Serialize(reqBody)
	if err != nil {
		return errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, sms.URL+"/mt/single", bytes.NewReader(b))
	if err != nil {
		return errors.E(op, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	token := "Bearer " + sms.Token
	req.Header.Add("Authorization", token)

	res, err := sms.client.Do(req)
	if err != nil {
		return errors.E(op, err)
	}
	defer res.Body.Close()

	resBody := &smsResponse{}

	if err := encoding.Deserialize(res.Body, resBody); err != nil {
		return errors.E(op, err)
	}

	if res.StatusCode > 299 {
		return errors.E(op, resBody.Message, errors.KindUnexpected)
	}

	return nil
}
func (sms *Backend) sendBulk(ctx context.Context, id string, recipients []string, message string) error {
	const op errors.Op = "backends/sms/Backend.sendBulk"

	reqBody := &bulkMSISDNRequest{
		MSISDN:    recipients,
		Message:   message,
		SenderID:  "PayPack",
		Reference: id,
	}

	b, err := encoding.Serialize(reqBody)
	if err != nil {
		return errors.E(op, err)
	}

	req, err := http.NewRequest(http.MethodPost, sms.URL+"/mt/bulk", bytes.NewReader(b))
	if err != nil {
		return errors.E(op, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	token := "Bearer " + sms.Token
	req.Header.Add("Authorization", token)

	res, err := sms.client.Do(req)
	if err != nil {
		return errors.E(op, err)
	}
	defer res.Body.Close()

	resBody := &smsResponse{}
	if err := encoding.Deserialize(res.Body, resBody); err != nil {
		return errors.E(op, err)
	}

	if res.StatusCode > 299 {
		return errors.E(op, resBody.Message, errors.KindUnexpected)
	}
	return nil
}

func (sms *Backend) authorize(appID, appSecret string) error {
	const op errors.Op = "backend/sms/backend.authorize"

	// assemble request
	body := &authRequest{AppID: appID, AppSecret: appSecret}

	bs, err := encoding.Serialize(body)
	if err != nil {
		return errors.E(op, err)
	}
	req, err := http.NewRequest(http.MethodPost, sms.URL+"/auth/", bytes.NewReader(bs))
	if err != nil {
		return errors.E(op, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// make request
	resp, err := sms.client.Do(req)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	// read response
	res := &authResponse{}

	if err := encoding.Deserialize(resp.Body, res); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	if err := res.Validate(); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	if resp.StatusCode > 299 {
		return errors.E(op, res.Message, errors.KindUnexpected)
	}

	sms.Token = res.AccessToken
	sms.Refresh = res.RefreshToken

	return nil
}

func (sms *Backend) refresh() error {
	const op errors.Op = "backend/sms/backend.refresh"

	// assemble request
	body := &refreshRequest{Token: sms.Refresh}

	bs, err := encoding.Serialize(body)
	if err != nil {
		return errors.E(op, err)
	}
	req, err := http.NewRequest(http.MethodPost, sms.URL+"/auth", bytes.NewReader(bs))
	if err != nil {
		return errors.E(op, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// make request
	resp, err := sms.client.Do(req)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	// read response
	res := &authResponse{}

	if err := encoding.Deserialize(resp.Body, res); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	if err := res.Validate(); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	if resp.StatusCode > 299 {
		return errors.E(op, res.Message, errors.KindUnexpected)
	}

	sms.Token = res.AccessToken
	sms.Refresh = res.RefreshToken

	return nil
}
