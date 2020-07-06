package payment_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/payment"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/payment/mocks"
	"github.com/rugwirobaker/paypack-backend/core/properties"
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

func newServer(svc payment.Service) *httptest.Server {
	mux := mux.NewRouter()
	var opts endpoints.HandlerOpts
	opts.Service = svc
	opts.Logger = log.NoOpLogger()
	endpoints.RegisterHandlers(mux, &opts)
	return httptest.NewServer(mux)
}

func TestPull(t *testing.T) {
	owners, owner := newOwnersStore()
	properties, property := newPropertiesStore(owner)
	invoices, invoice := newInvoiceStore(property)
	svc := newService(owners, properties, invoices)
	id := "123e4567-e89b-12d3-a456-000000000001"
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
			req: toJSON(payment.Payment{
				Code:   property.ID,
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
			req: toJSON(payment.Payment{
				Code:   property.ID,
				Phone:  "0785780891",
				Method: payment.MTN,
			}),
			status:      http.StatusBadRequest,
			contentType: contentType,
			res:         toJSON(map[string]string{"error": "amount must be greater than zero"}),
		},
		{
			desc: "missing house code",
			req: toJSON(payment.Payment{
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
			req: toJSON(payment.Payment{
				Code:   property.ID,
				Amount: invoice.Amount,
				Phone:  "0785780891",
			}),
			status:      http.StatusBadRequest,
			contentType: contentType,
			res:         toJSON(map[string]string{"error": "payment method must be specified"}),
		},
		{
			desc: "initialize payment with unsaved house code",
			req: toJSON(payment.Payment{
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
	owners, owner := newOwnersStore()
	properties, property := newPropertiesStore(owner)
	invoices, _ := newInvoiceStore(property)
	svc := newService(owners, properties, invoices)
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

func newService(ws owners.Repository, ps properties.Repository, vc invoices.Repository) payment.Service {
	var opts payment.Options
	opts.Owners = ws
	opts.Properties = ps
	opts.Invoices = vc
	opts.SMS = newSMSService()
	opts.Idp = mocks.NewIdentityProvider()
	opts.Backend = mocks.NewBackend()
	opts.Queue = mocks.NewQueue()
	opts.Transactions = mocks.NewTransactionsRepository()
	return payment.New(&opts)
}

func newPropertiesStore(owner owners.Owner) (properties.Repository, properties.Property) {
	var property properties.Property
	property.ID = uuid.New().ID()
	property.Due = 1000
	property.Owner = properties.Owner{ID: owner.ID}
	store := mocks.NewPropertyRepository()
	property, _ = store.Save(context.Background(), property)
	return store, property
}

func newOwnersStore() (owners.Repository, owners.Owner) {
	var owner owners.Owner
	owner.ID = uuid.New().ID()
	owner.Fname = "Jamie"
	owner.Lname = "Jones"
	owner.Phone = "0787205106"
	store := mocks.NewOwnersRepository()
	owner, _ = store.Save(context.Background(), owner)
	return store, owner
}
func newInvoiceStore(property properties.Property) (invoices.Repository, invoices.Invoice) {
	var invoice invoices.Invoice
	invoice.Status = invoices.Pending
	invoice.Amount = property.Due
	invoice.Property = property.ID
	creation := time.Now()
	var invs = map[string]invoices.Invoice{
		property.ID: {ID: 1, Amount: 1000, CreatedAt: creation, UpdatedAt: creation},
	}
	store := mocks.NewInvoiceRepository(invs)
	return store, invoice
}

func newSMSService() notifs.Service {
	var opts notifs.Options
	opts.IDP = uuid.New()
	opts.Backend = mocks.NewSMSBackend()
	opts.Store = mocks.NewSMSRepository()
	return notifs.New(&opts)
}
