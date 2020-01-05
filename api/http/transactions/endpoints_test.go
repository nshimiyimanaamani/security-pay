package transactions_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	email   = "user@gmail.com"
	token   = "token"
	wrong   = "wrong"
	wrongID = 0
)

var (
	contentType = "application/json"
	transaction = transactions.Transaction{Amount: 1000.00, Method: "BK", MadeFor: "1000-4433-3343", OwnerID: "1000-4433-3343"}
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	token       string
	body        io.Reader
}

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

func newService(tokens map[string]string) transactions.Service {
	repo := mocks.NewRepository()
	idp := mocks.NewIdentityProvider()
	opts := &transactions.Options{
		Repo: repo,
		Idp:  idp,
	}
	return transactions.New(opts)
}

func newServer(svc transactions.Service) *httptest.Server {
	mux := mux.NewRouter()
	opts := &endpoints.HandlerOpts{
		Service: svc,
		Logger:  log.NoOpLogger(),
	}
	endpoints.RegisterHandlers(mux, opts)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestRecord(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	id := "1"

	cases := []struct {
		desc        string
		req         string
		contentType string
		token       string
		status      int
	}{
		{
			desc:        "record a valid transaction",
			req:         toJSON(transactions.Transaction{ID: id, Amount: 100, Method: "bk", MadeFor: "1000", OwnerID: "1000"}),
			contentType: contentType,
			token:       token,
			status:      http.StatusCreated,
		},
		{
			desc:        "record transaction with invalid property",
			req:         toJSON(transactions.Transaction{Amount: 1000.00, Method: "BK"}),
			contentType: contentType,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "record transaction with invalid request format",
			req:         "{",
			contentType: contentType,
			token:       token,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "record transaction with empty JSON request",
			req:         "{}",
			contentType: contentType,
			token:       token,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "record transaction with empty request",
			req:         "",
			contentType: contentType,
			token:       token,
			status:      http.StatusBadRequest,
		},
		{
			desc:        "record transaction with missing content type",
			req:         toJSON(transactions.Transaction{ID: id, Amount: 100, Method: "bk", MadeFor: "1000", OwnerID: "1000"}),
			contentType: "",
			token:       token,
			status:      http.StatusUnsupportedMediaType,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/transactions", ts.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestRetrieve(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	ctx := context.Background()
	saved, err := svc.Record(ctx, transaction)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc   string
		id     string
		token  string
		status int
		res    string
	}{
		{
			desc:   "view existing transaction",
			id:     saved.ID,
			token:  token,
			status: http.StatusOK,
			res:    toJSON(saved),
		},
		{
			desc:   "view non-existent transaction",
			id:     strconv.FormatUint(wrongID, 10),
			token:  token,
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "transaction not found"}),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    fmt.Sprintf("%s/transactions/%s", ts.URL, tc.id),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}
}

func TestList(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	data := []transactions.Transaction{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		data = append(data, saved)
	}

	transactionURL := fmt.Sprintf("%s/transactions", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []transactions.Transaction
	}{
		{
			desc:   "get a list of transactions",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", transactionURL, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of transactions with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", transactionURL, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of transactions with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", transactionURL, 1, -5),
			res:    nil,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    tc.url,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var data transactions.TransactionPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Transactions, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Transactions))
	}
}

func TestListByProperty(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	data := []transactions.Transaction{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		data = append(data, saved)
	}

	transactionURL := fmt.Sprintf("%s/transactions", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []transactions.Transaction
	}{
		{
			desc:   "get a list of transactions",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?property=%s&offset=%d&limit=%d", transactionURL, transaction.MadeFor, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of transactions with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?property=%s&offset=%d&limit=%d", transactionURL, transaction.MadeFor, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of transactions with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?property=%s&offset=%d&limit=%d", transactionURL, transaction.MadeFor, 1, -5),
			res:    nil,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    tc.url,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var data transactions.TransactionPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code '%d' got '%d'", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Transactions, fmt.Sprintf("%s: expected body '%v' got '%v'", tc.desc, tc.res, data.Transactions))
	}
}

func TestListMethod(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	data := []transactions.Transaction{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		data = append(data, saved)
	}

	transactionURL := fmt.Sprintf("%s/transactions", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []transactions.Transaction
	}{
		{
			desc:   "get a list of transactions",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?method=%s&offset=%d&limit=%d", transactionURL, transaction.Method, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of transactions with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?method=%s&offset=%d&limit=%d", transactionURL, transaction.Method, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of transactions with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?method=%s&offset=%d&limit=%d", transactionURL, transaction.Method, 1, -5),
			res:    nil,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    tc.url,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var data transactions.TransactionPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Transactions, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Transactions))
	}
}

func TestMListByProperty(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	data := []transactions.Transaction{}

	for i := 0; i < 10; i++ {
		ctx := context.Background()
		saved, err := svc.Record(ctx, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		data = append(data, saved)
	}

	transactionURL := fmt.Sprintf("%s/mobile/transactions", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []transactions.Transaction
	}{
		{
			desc:   "get a list of transactions",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?property=%s", transactionURL, transaction.MadeFor),
			res:    data[0:10],
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    tc.url,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var data []transactions.Transaction
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data, fmt.Sprintf("%s: expected body '%v' got '%v'", tc.desc, tc.res, data))
	}
}
