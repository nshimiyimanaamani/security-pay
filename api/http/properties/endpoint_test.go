package properties_test

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
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	"github.com/rugwirobaker/paypack-backend/app/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/properties/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	wrongID     = 0
	email       = "user@example.com"
	token       = "token"
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

func newService() properties.Service {
	idp := mocks.NewIdentityProvider()
	repo := mocks.NewRepository()
	return properties.New(idp, repo)
}

func newServer(svc properties.Service) *httptest.Server {
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

func TestRegisterProperty(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: properties.Owner{ID: "1"},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := toJSON(property)
	invalidData := toJSON(properties.Property{})

	//fake id
	property.ID = "1"

	res := toJSON(property)
	invalidEntityRes := toJSON(Error{"invalid entity format"})
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
			desc:        "add a valid property",
			req:         data,
			token:       token,
			contentType: contentType,
			status:      http.StatusCreated,
			res:         res,
		},
		{
			desc:        "add property with invalid data",
			req:         invalidData,
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add property with invalid request format",
			req:         "{",
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add property with empty JSON request",
			req:         "{}",
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add property with empty request",
			req:         "",
			token:       token,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityRes,
		},
		{
			desc:        "add property with missing content type",
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
			url:         fmt.Sprintf("%s/properties", ts.URL),
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

func TestUpdateProperty(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: properties.Owner{ID: "1"},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	ctx := context.Background()
	saved, err := svc.RegisterProperty(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := Property{
		ID:    saved.ID,
		Owner: Owner{ID: saved.Owner.ID},
		Due:   saved.Due,
		Address: Address{
			Sector:  saved.Address.Sector,
			Cell:    saved.Address.Cell,
			Village: saved.Address.Village,
		},
	}

	data := toJSON(res)
	notFoundMessage := toJSON(Error{"non-existent property entity"})
	invalidEntityMessage := toJSON(Error{"invalid entity format"})
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
			desc:        "update existing property",
			req:         data,
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         data,
		},
		{
			desc:        "update non-existent property",
			req:         data,
			token:       token,
			id:          strconv.FormatUint(wrongID, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update property with invalid id",
			req:         data,
			token:       token,
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update property with invalid data format",
			req:         "{",
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update property with empty request",
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
			url:         fmt.Sprintf("%s/properties/%s", ts.URL, tc.id),
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

func TestRetrieveProperty(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: properties.Owner{ID: "1"},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	ctx := context.Background()
	saved, err := svc.RegisterProperty(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := Property{
		ID:    saved.ID,
		Owner: Owner{ID: saved.Owner.ID},
		Due:   saved.Due,
		Address: Address{
			Sector:  saved.Address.Sector,
			Cell:    saved.Address.Cell,
			Village: saved.Address.Village,
		},
	}

	data := toJSON(res)
	notFoundRes := toJSON(Error{"non-existent property entity"})

	cases := []struct {
		desc        string
		id          string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "view existing owner",
			id:     saved.ID,
			token:  token,
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "view non-existent owner",
			id:     strconv.FormatUint(wrongID, 10),
			token:  token,
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
		{
			desc:   "view property by passing invalid id",
			id:     "invalid",
			token:  token,
			status: http.StatusNotFound,
			res:    notFoundRes,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodGet,
			token:  tc.token,
			url:    fmt.Sprintf("%s/properties/%s", ts.URL, tc.id),
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

func TestListPropertiesByOwner(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{ID: uuid.New().ID()}

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := []Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := Property{
			ID:    saved.ID,
			Owner: Owner{ID: saved.Owner.ID},
			Due:   saved.Due,
			Address: Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []Property
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", transactionURL, owner.ID, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", transactionURL, owner.ID, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", transactionURL, owner.ID, 1, -5),
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
		var data Properties
		err = json.NewDecoder(res.Body).Decode(&data)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListPropertiesByCell(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	cell := "Gishushu"
	property := properties.Property{
		Owner: properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    cell,
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := []Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := Property{
			ID:    saved.ID,
			Owner: Owner{ID: saved.Owner.ID},
			Due:   saved.Due,
			Address: Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []Property
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?cell=%s&offset=%d&limit=%d", transactionURL, cell, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?cell=%s&offset=%d&limit=%d", transactionURL, cell, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?cell=%s&offset=%d&limit=%d", transactionURL, cell, 1, -5),
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
		var data Properties
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListPropertiesBySector(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	sector := "Remera"
	property := properties.Property{
		Owner: properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{
			Sector:  sector,
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := []Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := Property{
			ID:    saved.ID,
			Owner: Owner{ID: saved.Owner.ID},
			Due:   saved.Due,
			Address: Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []Property
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?sector=%s&offset=%d&limit=%d", transactionURL, sector, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?sector=%s&offset=%d&limit=%d", transactionURL, sector, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?sector=%s&offset=%d&limit=%d", transactionURL, sector, 1, -5),
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
		var data Properties
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListPropertiesByVillage(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	village := "Ingabo"

	property := properties.Property{
		Owner: properties.Owner{ID: uuid.New().ID()},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: village,
		},
		Due: float64(1000),
	}

	data := []Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.RegisterProperty(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := Property{
			ID:    saved.ID,
			Owner: Owner{ID: saved.Owner.ID},
			Due:   saved.Due,
			Address: Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []Property
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?village=%s&offset=%d&limit=%d", transactionURL, village, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?village=%s&offset=%d&limit=%d", transactionURL, village, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?village=%s&offset=%d&limit=%d", transactionURL, village, 1, -5),
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
		var data Properties
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

type Property struct {
	ID      string  `json:"id,omitempty"`
	Due     float64 `json:"due,string,omitempty"`
	Owner   Owner   `json:"owner,omitempty"`
	Address Address `json:"address,omitempty"`
}

type Address struct {
	Sector  string `json:"sector,omitempty"`
	Cell    string `json:"cell,omitempty"`
	Village string `json:"village,omitempty"`
}

type Owner struct {
	ID    string `json:"id,omitempty"`
	Fname string `json:"fname,omitempty"`
	Lname string `json:"lname,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type Properties struct {
	Properties []Property `json:"properties"`
	Total      uint64     `json:"total"`
	Offset     uint64     `json:"offset"`
	Limit      uint64     `json:"limit"`
}

type Error struct {
	Message string `json:"message"`
}
