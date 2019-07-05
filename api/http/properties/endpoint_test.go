package properties_test

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

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	wrongID     = 0
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

func newService() properties.Service {
	store := mocks.NewPropertyStore()
	idp := mocks.NewIdentityProvider()
	return properties.New(idp, store)
}

func newServer(svc properties.Service) *httptest.Server {
	mux := mux.NewRouter()
	endpoints.MakeEndpoint(mux)(svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestAddProperty(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: "Eugene Mugabo",
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
	}

	data := toJSON(property)
	invalidData := toJSON(properties.Property{Owner: "Gakwaya Daniel"})

	cases := []struct {
		desc        string
		req         string
		contentType string
		status      int
	}{
		{"add a valid property", data, contentType, http.StatusCreated},
		{"add property with invalid data", invalidData, contentType, http.StatusBadRequest},
		{"add property with invalid request format", "{", contentType, http.StatusBadRequest},
		{"add property with empty JSON request", "{}", contentType, http.StatusBadRequest},
		{"add property with empty request", "", contentType, http.StatusBadRequest},
		{"add property with missing content type", data, "", http.StatusUnsupportedMediaType},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		// body, _ := ioutil.ReadAll(res.Body)
		// data := strings.Trim(string(body), "\n")
		// assert.Equal(t, "", data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, "", data))
	}
}

func TestUpdateProperty(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: "Eugene Mugabo",
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
	}

	saved, err := svc.AddProperty(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := propRes{
		ID:      saved.ID,
		Owner:   saved.Owner,
		Sector:  saved.Sector,
		Cell:    saved.Cell,
		Village: saved.Village,
	}

	data := toJSON(res)
	notFoundMessage := toJSON(errRes{"non-existent entity"})
	invalidEntityMessage := toJSON(errRes{"invalid entity format"})
	unsupportedContentMessage := toJSON(errRes{"unsupported content type"})

	cases := []struct {
		desc        string
		req         string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "update existing property",
			req:         data,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         data,
		},
		{
			desc:        "update non-existent property",
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
			url:         fmt.Sprintf("%s/%s", ts.URL, tc.id),
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

func TestViewProperty(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: "Eugene Mugabo",
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
	}

	saved, err := svc.AddProperty(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := propRes{
		ID:      saved.ID,
		Owner:   saved.Owner,
		Sector:  saved.Sector,
		Cell:    saved.Cell,
		Village: saved.Village,
	}

	data := toJSON(res)
	notFoundMessage := toJSON(errRes{"non-existent entity"})

	cases := []struct {
		desc        string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "view existing transaction",
			id:     saved.ID,
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "view non-existent transaction",
			id:     strconv.FormatUint(wrongID, 10),
			status: http.StatusNotFound,
			res:    notFoundMessage,
		},
		{
			desc:   "view thing by passing invalid id",
			id:     "invalid",
			status: http.StatusNotFound,
			res:    notFoundMessage,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
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

// func TestListPropertiesByOwner(t *testing.T) {
// 	svc := newService()
// 	ts := newServer(svc)

// 	defer ts.Close()
// 	client := ts.Client()

// 	property := properties.Property{
// 		Owner: "Eugene Mugabo",
// 		Address: properties.Address{
// 			Sector:  "Remera",
// 			Cell:    "Gishushu",
// 			Village: "Ingabo",
// 		},
// 	}

// 	data := []propRes{}

// 	for i := 0; i < 100; i++ {
// 		saved, err := svc.AddProperty(property)
// 		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

// 		res := propRes{
// 			ID:      saved.ID,
// 			Owner:   saved.Owner,
// 			Sector:  saved.Sector,
// 			Cell:    saved.Cell,
// 			Village: saved.Village,
// 		}

// 		data = append(data, res)
// 	}

// 	transactionURL := fmt.Sprintf("%s/", ts.URL)

// 	cases := []struct {
// 		desc   string
// 		status int
// 		url    string
// 		res    []propRes
// 	}{
// 		{
// 			desc:   "get a list of transactions",
// 			status: http.StatusOK,
// 			url:    fmt.Sprintf("%s?offset=%d&limit=%d", transactionURL, 0, 5),
// 			res:    data[0:5],
// 		},
// 	}

// 	for _, tc := range cases {
// 		req := testRequest{
// 			client: client,
// 			method: http.MethodGet,
// 			url:    tc.url,
// 		}

// 		res, err := req.make()
// 		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
// 		var data propPageRes
// 		json.NewDecoder(res.Body).Decode(&data)
// 		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
// 		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
// 	}
// }

func TestListPropertiesBySector(t *testing.T) {}

func TestListPropertiesByCell(t *testing.T) {}

func TestListPropertiesByVillage(t *testing.T) {}

type propRes struct {
	ID      string `json:"id,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Sector  string `json:"sector,omitempty"`
	Cell    string `json:"cell,omitempty"`
	Village string `json:"village,omitempty"`
}

type propPageRes struct {
	Properties []propRes `json:"properties"`
	Total      uint64    `json:"total"`
	Offset     uint64    `json:"offset"`
	Limit      uint64    `json:"limit"`
}

type errRes struct {
	Message string `json:"message"`
}
