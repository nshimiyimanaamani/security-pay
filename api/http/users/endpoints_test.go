package users_test

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/stretchr/testify/require"

// 	//"errors"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/gorilla/mux"
// 	endpoints "github.com/nshimiyimanaamani/paypack-backend/api/http/users"
// 	"github.com/nshimiyimanaamani/paypack-backend/core/users"
// 	"github.com/nshimiyimanaamani/paypack-backend/core/users/mocks"
// 	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	contentType  = "application/json"
// 	invalidEmail = "userexample.com"
// )

// type testRequest struct {
// 	client      *http.Client
// 	method      string
// 	url         string
// 	contentType string
// 	token       string
// 	body        io.Reader
// }

// func (tr testRequest) make() (*http.Response, error) {
// 	req, err := http.NewRequest(tr.method, tr.url, tr.body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if tr.token != "" {
// 		req.Header.Set("Authorization", tr.token)
// 	}
// 	if tr.contentType != "" {
// 		req.Header.Set("Content-Type", tr.contentType)
// 	}
// 	return tr.client.Do(req)
// }

// func newService() users.Service {
// 	hasher := mocks.NewHasher()
// 	tempIdp := mocks.NewTempIdentityProvider()
// 	idp := mocks.NewIdentityProvider()
// 	store := mocks.NewUserStore()

// 	return users.New(hasher, tempIdp, idp, store)
// }

// func newServer(svc users.Service) *httptest.Server {
// 	mux := mux.NewRouter()
// 	opts := &endpoints.HandlerOpts{
// 		Service: svc,
// 		Logger:  log.NoOpLogger(),
// 	}
// 	endpoints.RegisterHandlers(mux, opts)
// 	return httptest.NewServer(mux)
// }

// func toJSON(data interface{}) string {
// 	jsonData, _ := json.Marshal(data)
// 	return string(jsonData)
// }

// func TestUserRegisterEndpoint(t *testing.T) {
// 	svc := newService()
// 	ts := newServer(svc)

// 	defer ts.Close()
// 	client := ts.Client()

// 	user := users.User{Username: "user@example.com", Password: "password", Cell: "cell", Village: "village", Sector: "sector"}

// 	data := toJSON(user)
// 	invalidData := toJSON(users.User{Username: invalidEmail, Password: "password"})
// 	invalidFieldData := fmt.Sprintf(`{"email": "%s", "pass": "%s"}`, user.Username, user.Password)

// 	res := toJSON(registrationRes{user.Username})
// 	conflictRes := toJSON(errRes{"user already exists"})
// 	invalidEntityRes := toJSON(errRes{"invalid entity format"})
// 	unsupportedContentRes := toJSON(errRes{"unsupported content type"})

// 	cases := []struct {
// 		desc        string
// 		req         string
// 		contentType string
// 		status      int
// 		res         string
// 	}{
// 		{"register new user", data, contentType, http.StatusCreated, res},
// 		{"register existing user", data, contentType, http.StatusConflict, conflictRes},
// 		{"register user with invalid email address", invalidData, contentType, http.StatusBadRequest, invalidEntityRes},
// 		{"register user with invalid request format", "{", contentType, http.StatusBadRequest, invalidEntityRes},
// 		{"register user with empty JSON request", "{}", contentType, http.StatusBadRequest, invalidEntityRes},
// 		{"register user with empty request", "", contentType, http.StatusBadRequest, invalidEntityRes},
// 		{"register user with invalid field name", invalidFieldData, contentType, http.StatusBadRequest, invalidEntityRes},
// 		{"register user with missing content type", data, "", http.StatusUnsupportedMediaType, unsupportedContentRes},
// 	}

// 	for _, tc := range cases {
// 		req := testRequest{
// 			client:      client,
// 			method:      http.MethodPost,
// 			url:         fmt.Sprintf("%s/users", ts.URL),
// 			contentType: tc.contentType,
// 			body:        strings.NewReader(tc.req),
// 		}

// 		res, err := req.make()
// 		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
// 		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
// 		body, err := ioutil.ReadAll(res.Body)
// 		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
// 		data := strings.Trim(string(body), "\n")
// 		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
// 		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
// 	}
// }

// func TestUserLoginEndpoint(t *testing.T) {
// 	svc := newService()
// 	ts := newServer(svc)

// 	defer ts.Close()
// 	client := ts.Client()

// 	user := users.User{Username: "user@example.com", Password: "password", Cell: "cell", Village: "village", Sector: "sector"}

// 	data := toJSON(user)
// 	invalidData := toJSON(users.User{Username: "user@example.com", Password: "invalid_password"})
// 	invalidEmailData := toJSON(users.User{Username: invalidEmail, Password: "password"})
// 	nonexistentData := toJSON(users.User{Username: "non-existentuser@example.com", Password: "pass"})

// 	tokenRes := toJSON(map[string]string{"token": user.Username})
// 	invalidEntityRes := toJSON(errRes{"invalid entity format"})
// 	invalidCredsRes := toJSON(errRes{"missing or invalid credentials provided"})
// 	unsupportedContentRes := toJSON(errRes{"unsupported content type"})

// 	_, err := svc.Register(user)
// 	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

// 	cases := []struct {
// 		desc        string
// 		req         string
// 		contentType string
// 		status      int
// 		res         string
// 	}{
// 		{
// 			desc:        "login with valid credentials",
// 			req:         data,
// 			contentType: contentType,
// 			status:      http.StatusCreated,
// 			res:         tokenRes,
// 		},
// 		{
// 			desc:        "login with invalid credentials",
// 			req:         invalidData,
// 			contentType: contentType,
// 			status:      http.StatusForbidden,
// 			res:         invalidCredsRes,
// 		},
// 		// technical debt
// 		{
// 			desc:        "login with invalid email address",
// 			req:         invalidEmailData,
// 			contentType: contentType,
// 			status:      http.StatusForbidden,
// 			res:         invalidCredsRes,
// 		},
// 		{
// 			desc:        "login non-existent user",
// 			req:         nonexistentData,
// 			contentType: contentType,
// 			status:      http.StatusForbidden,
// 			res:         invalidCredsRes,
// 		},
// 		{
// 			desc:        "login with invalid request format",
// 			req:         "{",
// 			contentType: contentType,
// 			status:      http.StatusBadRequest,
// 			res:         invalidEntityRes,
// 		},
// 		// {
// 		// 	desc:        "login with empty JSON request",
// 		// 	req:         "{}",
// 		// 	contentType: contentType,
// 		// 	status:      http.StatusBadRequest,
// 		// 	res:         invalidEntityRes,
// 		// },
// 		{
// 			desc:        "login with empty request",
// 			req:         "",
// 			contentType: contentType,
// 			status:      http.StatusBadRequest,
// 			res:         invalidEntityRes,
// 		},
// 		{
// 			desc:        "login with missing content type",
// 			req:         data,
// 			contentType: "",
// 			status:      http.StatusUnsupportedMediaType,
// 			res:         unsupportedContentRes,
// 		},
// 	}

// 	for _, tc := range cases {
// 		req := testRequest{
// 			client:      client,
// 			method:      http.MethodPost,
// 			url:         fmt.Sprintf("%s/users/tokens", ts.URL),
// 			contentType: tc.contentType,
// 			body:        strings.NewReader(tc.req),
// 		}

// 		res, err := req.make()
// 		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
// 		body, err := ioutil.ReadAll(res.Body)
// 		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
// 		token := strings.Trim(string(body), "\n")

// 		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
// 		assert.Equal(t, tc.res, token, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, token))
// 	}

// }

// type registrationRes struct {
// 	ID string `json:"id"`
// }

// type errRes struct {
// 	Message string `json:"message"`
// }
