package transactions_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/transactions/mocks"
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
	transaction = transactions.Transaction{Amount: "1000.00", Method: "BK", MadeFor: "1000-4433-3343", MadeBy: "1000-4433-3343"}
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
	auth := mocks.NewAuthBackend(tokens)
	idp := mocks.NewIdentityProvider()
	store := mocks.NewTransactionStore()
	return transactions.New(idp, store, auth)
}

func newServer(svc transactions.Service) *httptest.Server {
	mux := mux.NewRouter()
	endpoints.MakeAdapter(mux)(svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestRecordTransaction(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	data := toJSON(transaction)
	invalidData := toJSON(transactions.Transaction{Amount: "1000.00", Method: "BK"})

	cases := []struct {
		desc        string
		req         string
		contentType string
		token       string
		status      int
	}{
		{
			desc:        "record a valid transaction",
			req:         data,
			contentType: contentType,
			token:       token,
			status:      http.StatusCreated,
		},
		{
			desc:        "record transaction with invalid property",
			req:         invalidData,
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
			req:         data,
			contentType: "",
			token:       token,
			status:      http.StatusUnsupportedMediaType,
		},
		{
			desc:        "record a invalid transaction",
			req:         data,
			contentType: contentType,
			token:       wrong,
			status:      http.StatusForbidden,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/", ts.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestViewTransaction(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	strx, err := svc.RecordTransaction(token, transaction)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	trxRes := transRes{
		ID:       strx.ID,
		Property: strx.MadeFor,
		Owner:    strx.MadeBy,
		Amount:   strx.Amount,
		Method:   strx.Method,
	}

	data := toJSON(trxRes)
	notFoundRes := toJSON(errRes{"non-existent entity"})
	invalidCredsRes := toJSON(errRes{"missing or invalid credentials provided"})

	cases := []struct {
		desc   string
		id     string
		token  string
		status int
		res    string
	}{
		{
			desc:   "view existing transaction",
			id:     strx.ID,
			token:  token,
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "view non-existent transaction",
			id:     strconv.FormatUint(wrongID, 10),
			token:  token,
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
		{
			desc:   "view transaction with invalid token",
			id:     strx.ID,
			token:  wrong,
			status: http.StatusForbidden,
			res:    invalidCredsRes,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    fmt.Sprintf("%s/%s", ts.URL, tc.id),
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

func TestListTransactions(t *testing.T) {
	svc := newService(map[string]string{token: email})
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	data := []transRes{}

	for i := 0; i < 100; i++ {
		trx, err := svc.RecordTransaction(token, transaction)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		trxRes := transRes{
			ID:       trx.ID,
			Property: trx.MadeFor,
			Owner:    trx.MadeBy,
			Amount:   trx.Amount,
			Method:   trx.Method,
		}
		data = append(data, trxRes)
	}

	transactionURL := fmt.Sprintf("%s/", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []transRes
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
		{
			desc:   "list transactions with invalid credentials",
			token:  wrong,
			status: http.StatusForbidden,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", transactionURL, 0, 5),
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
		var data transPageRes
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Transactions, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Transactions))
	}
}

func TestListByProperty(t *testing.T) {}

func TestListByMethod(t *testing.T) {}

type errRes struct {
	Message string `json:"message"`
}

type transRes struct {
	ID           string            `json:"id,omitempty"`
	Property     string            `json:"property,omitempty"`
	Owner        string            `json:"owner,omitempty"`
	Amount       string            `json:"amount,omitempty"`
	Address      map[string]string `json:"address,omitempty"`
	Method       string            `json:"method,omitempty"`
	DateRecorded time.Time         `json:"recorded,omitempty"`
}

type transPageRes struct {
	Transactions []transRes `json:"transactions"`
	Total        uint64     `json:"total"`
	Offset       uint64     `json:"offset"`
	Limit        uint64     `json:"limit"`
}
