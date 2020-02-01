package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var user = auth.Credentials{Username: "jane", Role: "dev", Account: "paypack", Password: "password"}

func newService() auth.Service {

	repo := mocks.NewRepository(user)
	jwt := mocks.NewJWTProvider()
	hasher := mocks.NewHasher()
	opts := &auth.Options{Repo: repo, JWT: jwt, Hasher: hasher}
	return auth.New(opts)
}

func TestAuthenticate(t *testing.T) {
	svc := newService()

	token, err := svc.Login(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	r := mux.NewRouter()
	r.HandleFunc("/test", h)

	r.Use(Authenticate(svc))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	expected := http.StatusOK
	got := w.Result().StatusCode

	assert.Equal(t, expected, got, fmt.Sprintf("expected: '%d' got '%d'", expected, got))
}
