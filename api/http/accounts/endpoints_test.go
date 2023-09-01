package accounts_test

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
	endpoints "github.com/nshimiyimanaamani/paypack-backend/api/http/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/core/accounts/mocks"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	authmocks "github.com/nshimiyimanaamani/paypack-backend/core/auth/mocks"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/encrypt"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	wrongID     = 10
	wrongValue  = "wrong"
	token       = "rugwiro.account.dev"
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

func newAuthenticator() auth.Service {
	user := auth.Credentials{
		Username: "username",
		Password: "password",
		Role:     auth.Dev,
		Account:  "account",
	}
	opts := &auth.Options{
		Hasher:    authmocks.NewHasher(),
		Encrypter: encrypt.None(),
		Repo:      authmocks.NewRepository(user),
		JWT:       authmocks.NewJWTProvider(),
	}
	return auth.New(opts)
}

func newService() accounts.Service {
	repo := mocks.NewRepository()
	idp := mocks.NewIdentityProvider()
	opts := &accounts.Options{Repository: repo, IDP: idp}
	return accounts.New(opts)
}

func newServer(svc accounts.Service) *httptest.Server {
	mux := mux.NewRouter()
	opts := &endpoints.HandlerOpts{
		Service:       svc,
		Authenticator: newAuthenticator(),
		Logger:        log.NoOpLogger(),
	}
	endpoints.RegisterHandlers(mux, opts)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreate(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	id := "1"

	cases := []struct {
		desc        string
		req         string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "create valid account",
			token:       token,
			req:         toJSON(accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			contentType: contentType,
			status:      http.StatusCreated,
			res:         toJSON(accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
		},
		{
			desc:        " create account with missing account name",
			token:       token,
			req:         toJSON(accounts.Account{NumberOfSeats: 10, Type: accounts.Devs}),
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid account: missing name"}),
		},

		{
			desc:        "create account with missing type",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10}),
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid account: missing type"}),
		},
		{
			desc:        "record message with empty request",
			token:       token,
			req:         "",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "create account with missing content type",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         toJSON(map[string]string{"error": "invalid request: invalid content type"}),
		},
		{
			desc:        "create account with invalid token",
			req:         toJSON(accounts.Account{ID: id, Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/accounts", srv.URL),
			contentType: tc.contentType,
			token:       tc.token,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		//body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		//data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		//assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}
}

func TestUpdate(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	ctx := context.Background()

	saved, err := svc.Create(ctx, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc        string
		req         string
		id          string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "update valid account",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         toJSON(map[string]string{"message": fmt.Sprintf("account[%s]: updated", saved.ID)}),
		},
		{
			desc:        "update non-existant account",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			id:          strconv.FormatUint(wrongID, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         toJSON(map[string]string{"error": "account not found"}),
		},
		{
			desc:        "update account with invalid id",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         toJSON(map[string]string{"error": "account not found"}),
		},
		{
			desc:        " create account with missing account name",
			token:       token,
			req:         toJSON(accounts.Account{NumberOfSeats: 10, Type: accounts.Devs}),
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid account: missing name"}),
		},

		{
			desc:        "update account with missing type",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10}),
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid account: missing type"}),
		},
		{
			desc:        "update account with empty request",
			token:       token,
			req:         "",
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "update property with missing content type",
			token:       token,
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			id:          saved.ID,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         toJSON(map[string]string{"error": "invalid request: invalid content type"}),
		},
		{
			desc:        "update account with invalid token",
			req:         toJSON(accounts.Account{Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}),
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPut,
			url:         fmt.Sprintf("%s/accounts/%s", srv.URL, tc.id),
			contentType: tc.contentType,
			token:       tc.token,
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

func TestRetrieve(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	account := accounts.Account{ID: "gasabo.remera", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}
	ctx := context.Background()

	saved, err := svc.Create(ctx, account)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc   string
		id     string
		token  string
		status int
		res    string
	}{
		{
			desc:   "retrieve existant account",
			token:  token,
			id:     saved.ID,
			status: http.StatusOK,
			res:    toJSON(saved),
		},

		{
			desc:   "retrieve account with invalid id",
			token:  token,
			id:     "invalid",
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "account not found"}),
		},

		{
			desc:   "retrieve non-existent message",
			token:  token,
			id:     strconv.FormatUint(wrongID, 10),
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "account not found"}),
		},
		{
			desc:   "retrieve account with invalid token",
			id:     saved.ID,
			status: http.StatusUnauthorized,
			res:    toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    fmt.Sprintf("%s/accounts/%s", srv.URL, tc.id),
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

func TestList(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	account := accounts.Account{Name: "developers", NumberOfSeats: 10, Type: accounts.Devs}

	data := []accounts.Account{}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		account.ID = fmt.Sprintf("paypack.%d", i)
		saved, err := svc.Create(ctx, account)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
		data = append(data, saved)
	}

	accountsURL := fmt.Sprintf("%s/accounts", srv.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []accounts.Account
	}{
		{
			desc:   "get a list of accounts",
			token:  token,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", accountsURL, 0, 5),
			status: http.StatusOK,
			res:    data[0:5],
		},
		{
			desc:   "get a list of accounts with negative offset",
			token:  token,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", accountsURL, -1, 5),
			status: http.StatusBadRequest,
			res:    nil,
		},
		{
			desc:   "get a list of accounts with negative limit",
			token:  token,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", accountsURL, 1, -5),
			status: http.StatusBadRequest,
			res:    nil,
		},
		{
			desc:   "get a list of accounts with invalid token",
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", accountsURL, 0, 5),
			status: http.StatusUnauthorized,
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
		var data accounts.AccountPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Accounts, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Accounts))
	}
}

func TestDeactivate(t *testing.T) {}
