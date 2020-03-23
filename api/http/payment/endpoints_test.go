package payment_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/payment"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/payment/mocks"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
)

const (
	contentType = "application/json"
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

var idprovider = nanoid.New(&nanoid.Config{
	Alphabet: nanoid.Alphabet,
	Length:   nanoid.Length,
})

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}
	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}
	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}
	return tr.client.Do(req)
}

func newService(invoice payment.Invoice, properties []string) payment.Service {
	idp := mocks.NewIdentityProvider()
	backend := mocks.NewBackend()
	queue := mocks.NewQueue()
	repo := mocks.NewRepository(invoice, properties)
	opts := &payment.Options{Idp: idp, Backend: backend, Queue: queue, Repo: repo}
	return payment.New(opts)
}

func newServer(svc payment.Service) *httptest.Server {
	mux := mux.NewRouter()
	opts := &endpoints.HandlerOpts{
		Service: svc,
		Logger:  log.NoOpLogger(),
	}
	endpoints.RegisterHandlers(mux, opts)
	return httptest.NewServer(mux)
}

func TestInitialize(t *testing.T) {
	code := idprovider.ID()
	id := "123e4567-e89b-12d3-a456-000000000001"
	invoice := payment.Invoice{
		ID:     uint64(1000),
		Amount: float64(1000),
	}
	properties := []string{code}
	svc := newService(invoice, properties)
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	cases := []struct {
		desc        string
		req         string
		contentType string
		status      int
		res         string
	}{
		{
			desc: "initialize valid transaction",
			req: toJSON(payment.Transaction{
				Code:   code,
				Amount: invoice.Amount,
				Phone:  "0785780891",
				Method: payment.MTN,
			}),
			status:      http.StatusOK,
			contentType: contentType,
			res:         toJSON(map[string]string{"status": "success", "transaction_id": id, "transaction_state": "processing"}),
		},
		{
			desc: "empty transaction amount",
			req: toJSON(payment.Transaction{
				Code:   code,
				Phone:  "0785780891",
				Method: payment.MTN,
			}),
			status:      http.StatusBadRequest,
			contentType: contentType,
			res:         toJSON(map[string]string{"error": "amount must be greater than zero"}),
		},
		{
			desc: "missing house code",
			req: toJSON(payment.Transaction{
				Amount: invoice.Amount,
				Phone:  "0785780891",
				Method: payment.MTN,
			}),
			status:      http.StatusBadRequest,
			contentType: contentType,
			res:         toJSON(map[string]string{"error": "missing house code"}),
		},

		{
			desc: "missing payment method",
			req: toJSON(payment.Transaction{
				Code:   code,
				Amount: invoice.Amount,
				Phone:  "0785780891",
			}),
			status:      http.StatusBadRequest,
			contentType: contentType,
			res:         toJSON(map[string]string{"error": "payment method must be specified"}),
		},
		{
			desc: "initialize payment with unsaved house code",
			req: toJSON(payment.Transaction{
				Code:   idprovider.ID(),
				Amount: invoice.Amount,
				Phone:  "0785780891",
				Method: payment.MTN,
			}),
			status:      http.StatusNotFound,
			contentType: contentType,
			res:         toJSON(map[string]string{"error": "property not found"}),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/payment/initialize", srv.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}

}

func TestConfirm(t *testing.T) {
	code := uuid.New().ID()
	invoice := payment.Invoice{
		ID:     uint64(1000),
		Amount: float64(1000),
	}
	properties := []string{code}
	svc := newService(invoice, properties)
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	cases := []struct {
		desc        string
		req         string
		contentType string
		status      int
		res         string
	}{}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/owners", srv.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
