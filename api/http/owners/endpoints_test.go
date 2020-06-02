package owners_test

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
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/owners"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	authmocks "github.com/rugwirobaker/paypack-backend/core/auth/mocks"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/owners/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/encrypt"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	token       = "rugwiro.account.dev"
	wrong       = 0
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

func newService() owners.Service {
	opts := &owners.Options{
		Idp:  mocks.NewIdentityProvider(),
		Repo: mocks.NewRepository(),
	}
	return owners.New(opts)
}

func newServer(svc owners.Service) *httptest.Server {
	mux := mux.NewRouter()
	opts := &endpoints.HandlerOpts{
		Service:       svc,
		Logger:        log.NoOpLogger(),
		Authenticator: newAuthenticator(),
	}
	endpoints.RegisterHandlers(mux, opts)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestRegister(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}
	data := toJSON(owner)

	res := toJSON(Owner{ID: "1"})
	invalidData := toJSON(owners.Owner{})

	invalidEntityRes := toJSON(Error{"invalid owner entity"})
	unsupportedContentRes := toJSON(Error{"unsupported content type"})

	cases := []struct {
		desc        string
		req         string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "record valid owner",
			req:         data,
			token:       token,
			contentType: contentType,
			status:      http.StatusCreated,
			res:         res,
		},
		{
			desc:        "add owner with invalid data",
			req:         invalidData,
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add owner with invalid request format",
			req:         "{",
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add owner with empty JSON request",
			req:         "{}",
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add owner with empty request",
			req:         "",
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add owner with missing content type",
			req:         data,
			token:       token,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         unsupportedContentRes,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			token:       tc.token,
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

func TestRetrieve(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}

	ctx := context.Background()

	saved, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	data := toJSON(Owner{
		Fname: saved.Fname,
		Lname: saved.Lname,
		Phone: saved.Phone,
	})

	notFoundMessage := toJSON(Error{"owner entity not found"})

	cases := []struct {
		desc        string
		req         string
		token       string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "retrieve existing owner",
			req:         data,
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         toJSON(saved),
		},
		{
			desc:        "retrieve non-existent owner",
			req:         data,
			token:       token,
			id:          strconv.FormatUint(wrong, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "retrieve owner with invalid id",
			req:         data,
			token:       token,
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodGet,
			token:       tc.token,
			url:         fmt.Sprintf("%s/owners/%s", srv.URL, tc.id),
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

func TestUpdate(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}

	ctx := context.Background()

	saved, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	data := toJSON(Owner{
		Fname:     saved.Fname,
		Lname:     saved.Lname,
		Phone:     saved.Phone,
		Namespace: saved.Namespace,
	})

	notFoundMessage := toJSON(Error{"owner entity not found"})
	invalidEntityMessage := toJSON(Error{"invalid owner entity"})
	unsupportedContentMessage := toJSON(Error{"unsupported content type"})

	cases := []struct {
		desc        string
		req         string
		token       string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "update existing owner",
			req:         data,
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         toJSON(saved),
		},
		{
			desc:        "update non-existent owner",
			req:         data,
			token:       token,
			id:          strconv.FormatUint(wrong, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update owner with invalid id",
			req:         data,
			token:       token,
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update owner with invalid data format",
			req:         "{",
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update owner with empty request",
			req:         "",
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update thing without content type",
			req:         data,
			token:       token,
			id:          saved.ID,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         unsupportedContentMessage,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPut,
			token:       tc.token,
			url:         fmt.Sprintf("%s/owners/%s", srv.URL, tc.id),
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

func TestList(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}

	data := []Owner{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()

		saved, err := svc.Register(ctx, owner)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

		ow := Owner{
			ID:        saved.ID,
			Fname:     saved.Fname,
			Lname:     saved.Lname,
			Phone:     saved.Phone,
			Namespace: saved.Namespace,
		}

		data = append(data, ow)
	}

	transactionURL := fmt.Sprintf("%s/owners", srv.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []Owner
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?&offset=%d&limit=%d", transactionURL, 0, 5),
			res:    data[0:5],
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
		var data OwnersPage
		err = json.NewDecoder(res.Body).Decode(&data)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Owners, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Owners))
	}
}

func TestSearch(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}

	ctx := context.Background()
	saved, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	res := Owner{
		ID:        saved.ID,
		Fname:     saved.Fname,
		Lname:     saved.Lname,
		Phone:     saved.Phone,
		Namespace: saved.Namespace,
	}

	data := toJSON(res)
	notFoundMessage := toJSON(Error{"owner entity not found"})

	ownersURL := fmt.Sprintf("%s/owners/search", ts.URL)

	cases := []struct {
		desc   string
		token  string
		url    string
		status int
		res    string
	}{
		{
			desc:   "search existing owner",
			token:  token,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, res.Fname, res.Lname, res.Phone),
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "search owner with wrong first name",
			token:  token,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, "wrong", res.Lname, res.Phone),
			status: http.StatusNotFound,
			res:    notFoundMessage,
		},
		{
			desc:   "search owner with wrong lname name",
			token:  token,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, res.Fname, "wrong", res.Phone),
			status: http.StatusNotFound,
			res:    notFoundMessage,
		},
		{
			desc:   "search owner with wrong phone number",
			token:  token,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, res.Fname, res.Lname, "wrong"),
			status: http.StatusNotFound,
			res:    notFoundMessage,
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
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}
}

func TestRetrieveByPhone(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := owners.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882", Namespace: "namespace"}

	ctx := context.Background()
	saved, err := svc.Register(ctx, owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	res := Owner{
		ID:        saved.ID,
		Fname:     saved.Fname,
		Lname:     saved.Lname,
		Phone:     saved.Phone,
		Namespace: saved.Namespace,
	}

	data := toJSON(res)
	notFoundMessage := toJSON(Error{"owner entity not found"})

	ownersURL := fmt.Sprintf("%s/owners", ts.URL)

	cases := []struct {
		desc   string
		token  string
		url    string
		status int
		res    string
	}{
		{
			desc:   "search existing owner",
			token:  token,
			url:    fmt.Sprintf("%s?phone=%s", ownersURL, res.Phone),
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "search owner with wrong phone number",
			token:  token,
			url:    fmt.Sprintf("%s?&phone=%s", ownersURL, "wrong_number"),
			status: http.StatusNotFound,
			res:    notFoundMessage,
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
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}
}

type Owner struct {
	ID        string `json:"id,omitempty"`
	Fname     string `json:"fname,omitempty"`
	Lname     string `json:"lname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type OwnersPage struct {
	Owners       []Owner `json:"owners"`
	PageMetadata `json:",meta"`
}

type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

type Error struct {
	Message string `json:"message"`
}
