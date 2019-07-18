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
	idp := mocks.NewIdentityProvider()
	propStore := mocks.NewPropertyStore()
	ownerStore := mocks.NewOwnerStore()
	return properties.New(idp, ownerStore, propStore)
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
		Due: float64(1000),
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
		Due: float64(1000),
	}

	saved, err := svc.AddProperty(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := propRes{
		ID:      saved.ID,
		Owner:   saved.Owner,
		Due:     saved.Due,
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
		Due: float64(1000),
	}

	saved, err := svc.AddProperty(property)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := propRes{
		ID:      saved.ID,
		Owner:   saved.Owner,
		Due:     saved.Due,
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
			desc:   "view existing owner",
			id:     saved.ID,
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "view non-existent owner",
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

func TestListPropertiesByOwner(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	oid, err := svc.CreateOwner(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	property := properties.Property{
		Owner: oid,
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := []propRes{}

	for i := 0; i < 100; i++ {
		saved, err := svc.AddProperty(property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := propRes{
			ID:      saved.ID,
			Owner:   saved.Owner,
			Due:     saved.Due,
			Sector:  saved.Sector,
			Cell:    saved.Cell,
			Village: saved.Village,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/owners/properties", ts.URL)

	cases := []struct {
		desc   string
		status int
		url    string
		res    []propRes
	}{
		{
			desc:   "get a list of properties",
			status: http.StatusOK,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, oid, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, oid, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, oid, 1, -5),
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
		var data propPageRes
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
		Owner: "Gashuga John",
		Address: properties.Address{
			Sector:  sector,
			Cell:    "Gishushu",
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := []propRes{}

	for i := 0; i < 100; i++ {
		saved, err := svc.AddProperty(property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := propRes{
			ID:      saved.ID,
			Owner:   saved.Owner,
			Due:     saved.Due,
			Sector:  saved.Sector,
			Cell:    saved.Cell,
			Village: saved.Village,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/sectors", ts.URL)

	cases := []struct {
		desc   string
		status int
		url    string
		res    []propRes
	}{
		{
			desc:   "get a list of properties",
			status: http.StatusOK,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, sector, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, sector, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, sector, 1, -5),
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
		var data propPageRes
		json.NewDecoder(res.Body).Decode(&data)
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
		Owner: "Gashuga John",
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    cell,
			Village: "Ingabo",
		},
		Due: float64(1000),
	}

	data := []propRes{}

	for i := 0; i < 100; i++ {
		saved, err := svc.AddProperty(property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := propRes{
			ID:      saved.ID,
			Owner:   saved.Owner,
			Due:     saved.Due,
			Sector:  saved.Sector,
			Cell:    saved.Cell,
			Village: saved.Village,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/cells", ts.URL)

	cases := []struct {
		desc   string
		status int
		url    string
		res    []propRes
	}{
		{
			desc:   "get a list of properties",
			status: http.StatusOK,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, cell, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, cell, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, cell, 1, -5),
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
		var data propPageRes
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
		Owner: "Gashuga John",
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Gishushu",
			Village: village,
		},
		Due: float64(1000),
	}

	data := []propRes{}

	for i := 0; i < 100; i++ {
		saved, err := svc.AddProperty(property)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := propRes{
			ID:      saved.ID,
			Owner:   saved.Owner,
			Due:     saved.Due,
			Sector:  saved.Sector,
			Cell:    saved.Cell,
			Village: saved.Village,
		}

		data = append(data, res)
	}

	transactionURL := fmt.Sprintf("%s/villages", ts.URL)

	cases := []struct {
		desc   string
		status int
		url    string
		res    []propRes
	}{
		{
			desc:   "get a list of properties",
			status: http.StatusOK,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, village, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of properties with negative offset",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, village, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of properties with negative limit",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s/%s?offset=%d&limit=%d", transactionURL, village, 1, -5),
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
		var data propPageRes
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Properties, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Properties))
	}
}

func TestCreateOwner(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	invalidOwner := properties.Owner{}

	data := toJSON(owner)
	invalidData := toJSON(invalidOwner)

	cases := []struct {
		desc        string
		req         string
		contentType string
		status      int
	}{
		{"create a valid owner", data, contentType, http.StatusCreated},
		{"create owner with invalid property", invalidData, contentType, http.StatusBadRequest},
		{"create owner with invalid request format", "{", contentType, http.StatusBadRequest},
		{"create owner with empty JSON request", "{}", contentType, http.StatusBadRequest},
		{"create owner with empty request", "", contentType, http.StatusBadRequest},
		{"create owner with missing content type", data, "", http.StatusUnsupportedMediaType},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPost,
			url:         fmt.Sprintf("%s/owners/", ts.URL),
			contentType: tc.contentType,
			body:        strings.NewReader(tc.req),
		}

		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestUpdateOwner(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}
	invalidOwner := properties.Owner{}

	id, err := svc.CreateOwner(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := ownerRes{
		ID:    id,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}

	data := toJSON(res)
	invalidData := toJSON(invalidOwner)

	notFoundMessage := toJSON(errRes{"non-existent entity"})
	invalidEntityMessage := toJSON(errRes{"invalid entity format"})
	unsupportedContentMessage := toJSON(errRes{"unsupported content type"})

	ownersURL := fmt.Sprintf("%s/owners", ts.URL)

	cases := []struct {
		desc        string
		req         string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:        "update existing owner",
			req:         data,
			id:          id,
			contentType: contentType,
			status:      http.StatusOK,
			res:         data,
		},
		{
			desc:        "update non-existent owner",
			req:         data,
			id:          strconv.FormatUint(wrongID, 10),
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update owner with invalid id",
			req:         data,
			id:          "invalid",
			contentType: contentType,
			status:      http.StatusNotFound,
			res:         notFoundMessage,
		},
		{
			desc:        "update owner with invalid data format",
			req:         invalidData,
			id:          id,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update owner with invalid data format",
			req:         "{",
			id:          id,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update owner with empty request",
			req:         "",
			id:          id,
			contentType: contentType,
			status:      http.StatusBadRequest,
			res:         invalidEntityMessage,
		},
		{
			desc:        "update thing without content type",
			req:         data,
			id:          id,
			contentType: "",
			status:      http.StatusUnsupportedMediaType,
			res:         unsupportedContentMessage,
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      http.MethodPut,
			url:         fmt.Sprintf("%s/%s", ownersURL, tc.id),
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
func TestViewOwner(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}

	id, err := svc.CreateOwner(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := ownerRes{
		ID:    id,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}

	data := toJSON(res)
	notFoundMessage := toJSON(errRes{"non-existent entity"})

	ownersURL := fmt.Sprintf("%s/owners", ts.URL)

	cases := []struct {
		desc        string
		id          string
		contentType string
		status      int
		res         string
	}{
		{
			desc:   "view existing owner",
			id:     res.ID,
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "view non-existent owner",
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
			url:    fmt.Sprintf("%s/%s", ownersURL, tc.id),
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

func TestListOwners(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}

	data := []ownerRes{}

	for i := 0; i < 100; i++ {
		id, err := svc.CreateOwner(owner)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

		res := ownerRes{
			ID:    id,
			Fname: owner.Fname,
			Lname: owner.Lname,
			Phone: owner.Phone,
		}

		data = append(data, res)
	}

	ownersURL := fmt.Sprintf("%s/owners/", ts.URL)

	cases := []struct {
		desc   string
		status int
		url    string
		res    []ownerRes
	}{
		{
			desc:   "get a list of owners",
			status: http.StatusOK,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", ownersURL, 0, 5),
			res:    data[0:5],
		},
		{
			desc:   "get a list of owners with negative offset",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", ownersURL, -1, 5),
			res:    nil,
		},
		{
			desc:   "get a list of owners with negative limit",
			status: http.StatusBadRequest,
			url:    fmt.Sprintf("%s?offset=%d&limit=%d", ownersURL, 1, -5),
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
		var data ownerPageRes
		json.NewDecoder(res.Body).Decode(&data)
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.ElementsMatch(t, tc.res, data.Owners, fmt.Sprintf("%s: expected body %v got %v", tc.desc, tc.res, data.Owners))
	}
}

func TestSearchOwner(t *testing.T) {
	svc := newService()
	ts := newServer(svc)

	defer ts.Close()
	client := ts.Client()

	owner := properties.Owner{Fname: "James", Lname: "Torredo", Phone: "0784677882"}

	id, err := svc.CreateOwner(owner)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	res := ownerRes{
		ID:    id,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}

	data := toJSON(res)
	notFoundMessage := toJSON(errRes{"non-existent entity"})

	ownersURL := fmt.Sprintf("%s/owners/search/", ts.URL)

	cases := []struct {
		desc   string
		fname  string
		lname  string
		phone  string
		url    string
		status int
		res    string
	}{
		{
			desc:   "search existing owner",
			fname:  res.Fname,
			lname:  res.Lname,
			phone:  res.Phone,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, res.Fname, res.Lname, res.Phone),
			status: http.StatusOK,
			res:    data,
		},
		{
			desc:   "search owner with wrong first name",
			fname:  "wrong",
			lname:  res.Lname,
			phone:  res.Phone,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, "wrong", res.Lname, res.Phone),
			status: http.StatusNotFound,
			res:    notFoundMessage,
		},
		{
			desc:   "search owner with wrong lname name",
			fname:  res.Fname,
			lname:  "wrong",
			phone:  res.Phone,
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, res.Fname, "wrong", res.Phone),
			status: http.StatusNotFound,
			res:    notFoundMessage,
		},
		{
			desc:   "search owner with wrong phone number",
			fname:  res.Fname,
			lname:  res.Lname,
			phone:  "wrong",
			url:    fmt.Sprintf("%s?fname=%s&lname=%s&phone=%s", ownersURL, res.Fname, res.Lname, "wrong"),
			status: http.StatusNotFound,
			res:    notFoundMessage,
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
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		data := strings.Trim(string(body), "\n")
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
		assert.Equal(t, tc.res, data, fmt.Sprintf("%s: expected body %s got %s", tc.desc, tc.res, data))
	}

}

type propRes struct {
	ID      string  `json:"id,omitempty"`
	Owner   string  `json:"owner,omitempty"`
	Due     float64 `json:"due,string,omitempty"`
	Sector  string  `json:"sector,omitempty"`
	Cell    string  `json:"cell,omitempty"`
	Village string  `json:"village,omitempty"`
}

type ownerRes struct {
	ID    string `json:"id,omitempty"`
	Fname string `json:"fname,omitempty"`
	Lname string `json:"lname,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type propPageRes struct {
	Properties []propRes `json:"properties"`
	Total      uint64    `json:"total"`
	Offset     uint64    `json:"offset"`
	Limit      uint64    `json:"limit"`
}

type ownerPageRes struct {
	Owners []ownerRes `json:"owners"`
	Total  uint64     `json:"total"`
	Offset uint64     `json:"offset"`
	Limit  uint64     `json:"limit"`
}

type errRes struct {
	Message string `json:"message"`
}
