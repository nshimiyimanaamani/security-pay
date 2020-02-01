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

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/core/feedback"
	"github.com/rugwirobaker/paypack-backend/core/feedback/mocks"
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

	msg := feedback.Message{ID: "1", Title: "title", Body: "body", Creator: "0784677882"}

	validData := toJSON(msg)

	res := toJSON(feedback.Message{ID: "1", Title: "title", Body: "body", Creator: "0784677882"})

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
			req:         toJSON(feedback.Message{Title: "title"}),
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid message: missing body"}),
		},
		{
			desc:        "record message with invalid request format",
			req:         "{",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid message: missing title"}),
		},
		{
			desc:        "record message with empty JSON request",
			req:         "{}",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid message: missing title"}),
		},
		{
			desc:        "record message with empty request",
			req:         "",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "add message with missing content type",
			req:         validData,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         toJSON(map[string]string{"error": "invalid request: invalid content type"}),
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
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

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
			res:    toJSON(saved),
		},
		{
			desc:   "retrieve non-existent message",
			id:     strconv.FormatUint(wrongID, 10),
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "message not found"}),
		},
		{
			desc:   "retrieve message by passing invalid id",
			id:     "invalid",
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "message not found"}),
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
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

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
			res:    toJSON(map[string]string{"message": "message deleted"}),
		},
		{
			desc:   "delete non-existent message",
			id:     strconv.FormatUint(wrongID, 10),
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "message not found"}),
		},
		{
			desc:   "delete message by passing invalid id",
			id:     "invalid",
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "message not found"}),
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
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

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
			req:         toJSON(saved),
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         expected,
		},
		{
			desc:        "update non-existent message",
			req:         toJSON(saved),
			id:          strconv.FormatUint(wrongID, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         toJSON(map[string]string{"error": "message not found"}),
		},
		{
			desc:        "update message with invalid id",
			req:         toJSON(saved),
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         toJSON(map[string]string{"error": "message not found"}),
		},
		{
			desc:        "update message with invalid data format",
			req:         "{",
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid message: missing title"}),
		},
		{
			desc:        "update message with empty request",
			req:         "",
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "update thing without content type",
			req:         toJSON(saved),
			id:          saved.ID,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         toJSON(map[string]string{"error": "invalid request: invalid content type"}),
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

func TestList(t *testing.T) {
	svc := newService()
	srv := newServer(svc)

	defer srv.Close()
	client := srv.Client()

	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	data := []feedback.Message{}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		saved, err := svc.Record(ctx, &msg)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
		data = append(data, *saved)
	}

	messagesURL := fmt.Sprintf("%s/feedback", srv.URL)

	cases := []struct {
		desc   string
		status int
		url    string
		res    []feedback.Message
	}{
		{
			desc:   "get a list of messages",
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", messagesURL, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of messages with negative offset",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", messagesURL, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of messages with negative limit",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", messagesURL, 1, -5),
			res:    nil,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			url:    tc.url,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var data feedback.MessagePage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Messages, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Messages))
	}

}
