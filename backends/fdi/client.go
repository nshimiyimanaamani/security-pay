package fdi

import (
	"context"
	"time"

	"github.com/quarksgroup/paypack-go/paypack"
	"github.com/quarksgroup/paypack-go/paypack/api"
	"github.com/rugwirobaker/paypack-backend/core/payment"
)

// Timeout sets the default client timeout
const Timeout = 30 * time.Second

type Service struct {
	client      *api.Client
	token       *payment.PaymentTk
	WebhookMode string
}

func New(cli *api.Client, id, secret, mode string) (srv *Service, err error) {

	srv = &Service{
		client:      cli,
		WebhookMode: mode,
	}

	srv.token, err = srv.Login(context.Background(), id, secret)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Service) Login(ctx context.Context, client_id, client_secret string) (*payment.PaymentTk, error) {

	token, err := srv.client.Login(ctx, client_id, client_secret)

	if err != nil {
		return nil, err
	}
	return &payment.PaymentTk{
		Access:  token.Access,
		Refresh: token.Refresh,
		Expires: token.Expires,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refresh_token string) (*payment.PaymentTk, error) {

	in := &paypack.Token{
		Refresh: refresh_token,
	}

	token, err := s.client.Refresh(ctx, in)

	if err != nil {
		return nil, err
	}

	return &payment.PaymentTk{
		Access:  token.Access,
		Refresh: token.Refresh,
		Expires: token.Expires,
	}, nil
}

func (srv *Service) Push(ctx context.Context, txn *payment.TxRequest) (out *payment.TxResponse, err error) {

	if !time.Unix(srv.token.Expires, 0).After(time.Now().UTC().Add(-time.Minute)) {
		srv.token, err = srv.Refresh(ctx, srv.token.Refresh)
		if err != nil {
			return nil, err
		}
	}

	ctx = context.WithValue(ctx, paypack.TokenKey{}, &paypack.Token{
		Access: srv.token.Access,
	})

	tx := &paypack.TransactionRequest{
		Amount: txn.Amount,
		Number: txn.MSISDN,
		Mode:   srv.WebhookMode,
	}

	res, err := srv.client.Cashout(ctx, tx)

	if err != nil {
		return nil, err
	}

	out = &payment.TxResponse{
		TxID:    res.Ref,
		Status:  res.Status,
		Message: res.Status,
		TxState: payment.ToTxState[res.Status],
	}

	return out, err
}

func (srv *Service) Pull(ctx context.Context, txn *payment.TxRequest) (out *payment.TxResponse, err error) {

	if !time.Unix(srv.token.Expires, 0).After(time.Now().UTC().Add(-time.Minute)) {
		srv.token, err = srv.Refresh(ctx, srv.token.Refresh)
		if err != nil {
			return nil, err
		}
	}

	ctx = context.WithValue(ctx, paypack.TokenKey{}, &paypack.Token{
		Access: srv.token.Access,
	})

	tx := &paypack.TransactionRequest{
		Amount: txn.Amount,
		Number: txn.MSISDN,
		Mode:   srv.WebhookMode,
	}

	res, err := srv.client.Cashin(ctx, tx)

	if err != nil {
		return nil, err
	}

	out = &payment.TxResponse{
		TxID:    res.Ref,
		Status:  res.Status,
		Message: res.Status,
		TxState: payment.ToTxState[res.Status],
	}

	return out, err
}

var _ payment.Client = (*Service)(nil)
