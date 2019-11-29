package feedback_test

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
	"time"

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
	"github.com/rugwirobaker/paypack-backend/app/feedback/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	wrongID     = 10
	wrongValue  = "wrong"
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

func newService() feedback.Service {
	opts := &feedback.Options{
		Idp:  mocks.NewIdentityProvider(),
		Repo: mocks.NewRepository(),
	}
	return feedback.New(opts)
}

func newServer(svc feedback.Service) *httptest.Server {
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
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	validData := toJSON(msg)
	invalidData := toJSON(feedback.Message{Title: "title"})

	res := toJSON(msgRes{ID: "1"})
	invalidEntityRes := toJSON(errRes{"invalid message entity"})
	unsupportedContentRes := toJSON(errRes{"unsupported content type"})

	cases := []struct {
		desc        string
		req         string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "record valid message",
			req:         validData,
			contentType: contentType,
			status:      http.StatusCreated,
			res:         res,
		},
		{
			desc:        "record message with invalid data",
			req:         invalidData,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "record message with invalid request format",
			req:         "{",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "record message with empty JSON request",
			req:         "{}",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "record message with empty request",
			req:         "",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add property with missing content type",
			req:         validData,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         unsupportedContentRes,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/feedback", srv.URL),
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

	ctx := context.Background()

	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	saved, err := svc.Record(ctx, &msg)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	expected := msgRes{
		ID:        saved.ID,
		Title:     saved.Title,
		Body:      saved.Body,
		Creator:   saved.Creator,
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
	}

	data := toJSON(expected)
	notFoundRes := toJSON(errRes{"message entity not found"})

	cases := []struct {
		desc        string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "retrieve existing message",
			id:     saved.ID,
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "retrieve non-existent message",
			id:     strconv.FormatUint(wrongID, 10),
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
		{
			desc:   "retrieve message by passing invalid id",
			id:     "invalid",
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			url:    fmt.Sprintf("%s/feedback/%s", srv.URL, tc.id),
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

func TestDelete(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	ctx := context.Background()

	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	saved, err := svc.Record(ctx, &msg)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	notFoundRes := toJSON(errRes{"message entity not found"})

	expected := toJSON(map[string]string{"message": "message deleted"})

	cases := []struct {
		desc        string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "delete existing message",
			id:     saved.ID,
			status: http.StatusOK,
			res:    expected,
		},
		{
			desc:   "delete non-existent message",
			id:     strconv.FormatUint(wrongID, 10),
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
		{
			desc:   "delete message by passing invalid id",
			id:     "invalid",
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodDelete,
			url:    fmt.Sprintf("%s/feedback/%s", srv.URL, tc.id),
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

func TestUpdate(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	ctx := context.Background()

	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	saved, err := svc.Record(ctx, &msg)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := msgRes{
		ID:        saved.ID,
		Title:     saved.Title,
		Body:      saved.Body,
		Creator:   saved.Creator,
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
	}

	data := toJSON(res)

	notFoundMessage := toJSON(errRes{"message entity not found"})
	invalidEntityMessage := toJSON(errRes{"invalid message entity"})
	unsupportedContentMessage := toJSON(errRes{"unsupported content type"})

	expected := toJSON(map[string]string{"message": "message updated"})

	cases := []struct {
		desc        string
		req         string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "update existing message",
			req:         data,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         expected,
		},
		{
			desc:        "update non-existent message",
			req:         data,
			id:          strconv.FormatUint(wrongID, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update property with invalid id",
			req:         data,
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update property with invalid data format",
			req:         "{",
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update property with empty request",
			req:         "",
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update thing without content type",
			req:         data,
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
			url:         fmt.Sprintf("%s/feedback/%s", srv.URL, tc.id),
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

type msgRes struct {
	ID        string    `json:"id"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty"`
}

type errRes struct {
	Message string `json:"message"`
}
