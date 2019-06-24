package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"os"
	"strings"
	"testing"
	//"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/users/mocks"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/stretchr/testify/assert"
)


var (
	user = models.NewUser("user@example.com", "password")
	invalidEmail = "userexample.com"
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


func newService() users.Service {
	hasher:= mocks.NewHasher()
	tempIdp := mocks.NewTempIdentityProvider()
	idp := mocks.NewIdentityProvider()
	store:= mocks.NewUserStore()

	return users.New(hasher,tempIdp, idp, store)
}

func newServer(api *API,f func(http.ResponseWriter, *http.Request) error) *httptest.Server {
	mux:= api.handler(f)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestUserRegisterEndpoint(t *testing.T){
	svc := newService()
	api:= New(svc)
	ts := newServer(api, api.UserRegisterEndpoint)
	defer ts.Close()
	client := ts.Client()

	data := toJSON(user)
	invalidData := toJSON(models.User{Email: invalidEmail, Password: "password"})
	invalidFieldData := fmt.Sprintf(`{"email": "%s", "pass": "%s"}`, user.Email, user.Password)

	cases := []struct {
		desc        string
		req         string
		contentType string
		status      int
	}{
		{"register new user", data, contentType, http.StatusCreated},
		{"register existing user", data, contentType, http.StatusConflict},
		{"register user with invalid email address", invalidData, contentType, http.StatusBadRequest},
		{"register user with invalid request format", "{", contentType, http.StatusBadRequest},
		{"register user with empty JSON request", "{}", contentType, http.StatusBadRequest},
		{"register user with empty request", "", contentType, http.StatusBadRequest},
		{"register user with invalid field name", invalidFieldData, contentType, http.StatusBadRequest},
		{"register user with missing content type", data, "", http.StatusUnsupportedMediaType},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/users", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}



func TestUserLoginEndpoint(t *testing.T){
	svc := newService()
	api:= New(svc)
	ts := newServer(api, api.UserLoginEndpoint)
	defer ts.Close()
	client := ts.Client()

	tokenData := toJSON(map[string]string{"token": user.Email})
	data := toJSON(user)
	invalidData := toJSON(models.User{Email:"user@example.com", Password:"invalid_password"})
	invalidEmailData := toJSON(models.User{Email: invalidEmail, Password: "password"})
	nonexistentData := toJSON(models.User{Email:"non-existentuser@example.com", Password:"pass"})

	svc.Register(user)

	cases := []struct {
		desc        string
		req         string
		contentType string
		status      int
		res         string
	}{
		{"login with valid credentials", data, contentType, http.StatusCreated, tokenData},
		{"login with invalid credentials", invalidData, contentType, http.StatusForbidden, ""},
		{"login with invalid email address", invalidEmailData, contentType, http.StatusBadRequest, ""},
		{"login non-existent user", nonexistentData, contentType, http.StatusForbidden, ""},
		{"login with invalid request format", "{", contentType, http.StatusBadRequest, ""},
		{"login with empty JSON request", "{}", contentType, http.StatusBadRequest, ""},
		{"login with empty request", "", contentType, http.StatusBadRequest, ""},
		{"login with missing content type", data, "", http.StatusUnsupportedMediaType, ""},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/tokens", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		token := strings.Trim(string(body), "\n")

		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, token, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, token))
	}

}