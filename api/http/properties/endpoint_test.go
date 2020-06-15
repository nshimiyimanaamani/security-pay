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
	"time"

	"github.com/gorilla/mux"
	endpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	authmocks "github.com/rugwirobaker/paypack-backend/core/auth/mocks"
	"github.com/rugwirobaker/paypack-backend/core/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/properties/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/encrypt"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	wrongID     = 0
	email       = "user@example.com"
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

func newService(owners map[string]properties.Owner) properties.Service {
	idp := mocks.NewIdentityProvider()
	props := mocks.NewRepository(owners)
	return properties.New(idp, props)
}

func newServer(svc properties.Service) *httptest.Server {
	mux := mux.NewRouter()
	opts := &endpoints.HandlerOpts{
		Service:       svc,
		Logger:        log.NoOpLogger(),
		Authenticator: newAuthenticator(),
	}
	endpoints.RegisterHandlers(mux, opts)
	return httptest.NewServer(mux)
}

func makeOwners(owner properties.Owner) map[string]properties.Owner {
	owners := make(map[string]properties.Owner)
	owners[owner.ID] = owner
	return owners
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestRegister(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
	}

	data := toJSON(property)
	invalidData := toJSON(properties.Property{})

	//fake id
	property.ID = "1"

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
			token:       token,
			req:         data,
			contentType: contentType,
			status:      http.StatusCreated,
			res:         toJSON(property),
		},
		{
			desc:        "add property with invalid data",
			token:       token,
			req:         invalidData,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid property: missing owner"}),
		},
		{
			desc:        "add property with invalid request format",
			token:       token,
			req:         "{",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "add property with empty JSON request",
			token:       token,
			req:         "{}",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid property: missing owner"}),
		},
		{
			desc:        "add property with empty request",
			token:       token,
			req:         "",
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "add property with missing content type",
			token:       token,
			req:         toJSON(property),
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         toJSON(map[string]string{"error": "invalid request: invalid content type"}),
		},
		{
			desc:        "add property with empty token",
			req:         data,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},

		{
			desc:        "add property with invalid token",
			req:         data,
			token:       "invalid",
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         toJSON(map[string]string{"error": "access denied: invalid token"}),
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

func TestUpdate(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
	}

	ctx := context.Background()
	saved, err := svc.Register(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	res := properties.Property{
		ID:    saved.ID,
		Owner: owner,
		Due:   saved.Due,
		Address: properties.Address{
			Sector:  saved.Address.Sector,
			Cell:    saved.Address.Cell,
			Village: saved.Address.Village,
		},
		Namespace:  saved.Namespace,
		RecordedBy: saved.RecordedBy,
		Occupied:   saved.Occupied,
	}

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
			req:         toJSON(res),
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusOK,
			res:         toJSON(map[string]string{"message": fmt.Sprintf("property: '%s' successfully updated", res.ID)}),
		},
		{
			desc:        "update non-existent property",
			req:         toJSON(res),
			token:       token,
			id:          strconv.FormatUint(wrongID, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         toJSON(map[string]string{"error": "property not found"}),
		},
		{
			desc:        "update property with invalid id",
			req:         toJSON(res),
			token:       token,
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         toJSON(map[string]string{"error": "property not found"}),
		},
		{
			desc:        "update property with invalid data format",
			req:         "{",
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "update property with empty request",
			req:         "",
			token:       token,
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         toJSON(map[string]string{"error": "invalid request: wrong data format"}),
		},
		{
			desc:        "update thing without content type",
			req:         toJSON(res),
			token:       token,
			id:          saved.ID,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         toJSON(map[string]string{"error": "invalid request: invalid content type"}),
		},
		{
			desc:        "update property with empty token",
			req:         toJSON(res),
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},

		{
			desc:        "update property with invalid token",
			req:         toJSON(res),
			token:       "invalid",
			id:          saved.ID,
			contentType: contentType,
			status:      http.StatusUnauthorized,
			res:         toJSON(map[string]string{"error": "access denied: invalid token"}),
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

func TestRetrieve(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	now := time.Now()

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due:        float64(1000),
		Namespace:  "kigali.gasabo.remera",
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	ctx := context.Background()
	saved, err := svc.Register(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	res := properties.Property{
		ID:    saved.ID,
		Owner: properties.Owner{ID: saved.Owner.ID},
		Due:   saved.Due,
		Address: properties.Address{
			Sector:  saved.Address.Sector,
			Cell:    saved.Address.Cell,
			Village: saved.Address.Village,
		},
		Namespace:  saved.Namespace,
		RecordedBy: saved.RecordedBy,
		Occupied:   saved.Occupied,
		CreatedAt:  saved.CreatedAt,
		UpdatedAt:  saved.UpdatedAt,
	}

	cases := []struct {
		desc        string
		id          string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "view existing property",
			id:     saved.ID,
			token:  token,
			status: http.StatusOK,
			res:    toJSON(res),
		},
		{
			desc:   "view non-existent property",
			id:     strconv.FormatUint(wrongID, 10),
			token:  token,
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "property not found"}),
		},
		{
			desc:   "view property by passing invalid id",
			id:     "invalid",
			token:  token,
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "property not found"}),
		},
		{
			desc:   "view property with emtpy token",
			id:     saved.ID,
			status: http.StatusUnauthorized,
			res:    toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},
		{
			desc:   "view property with invalid token",
			id:     saved.ID,
			token:  "invalid",
			status: http.StatusUnauthorized,
			res:    toJSON(map[string]string{"error": "access denied: invalid token"}),
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

func TestDelete(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	now := time.Now()

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	ctx := context.Background()
	saved, err := svc.Register(ctx, property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc        string
		id          string
		token       string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "delete existing property",
			id:     saved.ID,
			token:  token,
			status: http.StatusOK,
			res:    toJSON(map[string]string{"message": "property deleted"}),
		},
		{
			desc:   "delete non-existent property",
			id:     strconv.FormatUint(wrongID, 10),
			token:  token,
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "property not found"}),
		},
		{
			desc:   "delete property by passing invalid id",
			id:     "invalid",
			token:  token,
			status: http.StatusNotFound,
			res:    toJSON(map[string]string{"error": "property not found"}),
		},
		{
			desc:   "delete property with empty token",
			id:     saved.ID,
			status: http.StatusUnauthorized,
			res:    toJSON(map[string]string{"error": "access denied: missing authorization token"}),
		},
		{
			desc:   "delete existing property",
			id:     saved.ID,
			token:  "invalid",
			status: http.StatusUnauthorized,
			res:    toJSON(map[string]string{"error": "access denied: invalid token"}),
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			method: http.MethodDelete,
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

func TestListByOwner(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()
	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
	}

	data := []properties.Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Register(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

		res := properties.Property{
			ID:    saved.ID,
			Owner: saved.Owner,
			Due:   saved.Due,
			Address: properties.Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
			Namespace:  saved.Namespace,
			RecordedBy: saved.RecordedBy,
			Occupied:   saved.Occupied,
		}

		data = append(data, res)
	}

	propertiesURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []properties.Property
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", propertiesURL, owner.ID, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", propertiesURL, owner.ID, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", propertiesURL, owner.ID, 1, -5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with empty token",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", propertiesURL, owner.ID, 0, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with invalid token",
			token:  "invalid",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?owner=%s&offset=%d&limit=%d", propertiesURL, owner.ID, 0, 5),
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
		var data properties.PropertyPage
		err = json.NewDecoder(res.Body).Decode(&data)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListByCell(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	cell := "Gishushu"
	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    cell,
			Village: "Ingabo",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
	}

	data := []properties.Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Register(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

		res := properties.Property{
			ID:    saved.ID,
			Owner: owner,
			Due:   saved.Due,
			Address: properties.Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
			Namespace:  saved.Namespace,
			RecordedBy: saved.RecordedBy,
			Occupied:   saved.Occupied,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []properties.Property
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
		{
			desc:   "get a list of properties with empty token",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?cell=%s&offset=%d&limit=%d", transactionURL, cell, 0, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with invalid token",
			token:  "invalid",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?cell=%s&offset=%d&limit=%d", transactionURL, cell, 0, 5),
			res:    nil,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client: client,
			token:  tc.token,
			method: http.MethodGet,
			url:    tc.url,
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		var data properties.PropertyPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListBySector(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}
	svc := newService(makeOwners(owner))

	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	sector := "Remera"
	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  sector,
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
	}

	data := []properties.Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Register(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

		res := properties.Property{
			ID:    saved.ID,
			Owner: owner,
			Due:   saved.Due,
			Address: properties.Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
			Namespace:  saved.Namespace,
			RecordedBy: saved.RecordedBy,
			Occupied:   saved.Occupied,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []properties.Property
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
		{
			desc:   "get a list of properties with empty token",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?sector=%s&offset=%d&limit=%d", transactionURL, sector, 0, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with invalid token",
			token:  "invalid",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?sector=%s&offset=%d&limit=%d", transactionURL, sector, 0, 5),
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
		var data properties.PropertyPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListByVillage(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}

	svc := newService(makeOwners(owner))
	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	village := "Ingabo"

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: village,
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: uuid.New().ID(),
		Occupied:   true,
	}

	data := []properties.Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Register(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

		res := properties.Property{
			ID:    saved.ID,
			Owner: owner,
			Due:   saved.Due,
			Address: properties.Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
			Namespace:  saved.Namespace,
			RecordedBy: saved.RecordedBy,
			Occupied:   saved.Occupied,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []properties.Property
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
		{
			desc:   "get a list of properties with empty token",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?village=%s&offset=%d&limit=%d", transactionURL, village, 0, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with invalid token",
			token:  "invalid",
			status: http.StatusUnauthorized,
			url:    fmt.Sprintf("%s?village=%s&offset=%d&limit=%d", transactionURL, village, 0, 5),
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
		var data properties.PropertyPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestListByRecorder(t *testing.T) {
	owner := properties.Owner{ID: uuid.New().ID()}

	svc := newService(makeOwners(owner))
	ts := newServer(svc)
	defer ts.Close()
	client := ts.Client()

	user := uuid.New().ID()

	property := properties.Property{
		Owner: owner,
		Address: properties.Address{
			Sector:  "remera",
			Cell:    "cell",
			Village: "village",
		},
		Namespace:  "kigali.gasabo.remera",
		Due:        float64(1000),
		RecordedBy: user,
		Occupied:   true,
	}

	data := []properties.Property{}

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		saved, err := svc.Register(ctx, property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

		res := properties.Property{
			ID:    saved.ID,
			Owner: owner,
			Due:   saved.Due,
			Address: properties.Address{
				Sector:  saved.Address.Sector,
				Cell:    saved.Address.Cell,
				Village: saved.Address.Village,
			},
			Namespace:  saved.Namespace,
			RecordedBy: saved.RecordedBy,
			Occupied:   saved.Occupied,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/properties", ts.URL)

	cases := []struct {
		desc   string
		token  string
		status int
		url    string
		res    []properties.Property
	}{
		{
			desc:   "get a list of properties",
			token:  token,
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?user=%s&offset=%d&limit=%d", transactionURL, user, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?user=%s&offset=%d&limit=%d", transactionURL, user, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			token:  token,
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?user=%s&offset=%d&limit=%d", transactionURL, user, 1, -5),
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
		var data properties.PropertyPage
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}
