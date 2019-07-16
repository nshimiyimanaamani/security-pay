package properties

import "net/http"

// Property defines a property(house) data model
type addPropertyRes struct {
	ID string `json:"id"`
}

func (res addPropertyRes) Code() int {
	return http.StatusCreated
}

func (res addPropertyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res addPropertyRes) Empty() bool {
	return false
}

type updatePropertyRes struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner,omitempty"`
	Due     float64 `json:"due,omitempty"`
	Sector  string  `json:"sector,omitempty"`
	Cell    string  `json:"cell,omitempty"`
	Village string  `json:"village,omitempty"`
}

func (res updatePropertyRes) Code() int {
	return http.StatusOK
}

func (res updatePropertyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res updatePropertyRes) Empty() bool {
	return false
}

type viewPropRes struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner,omitempty"`
	Due     float64 `json:"due,omitempty"`
	Sector  string  `json:"sector,omitempty"`
	Cell    string  `json:"cell,omitempty"`
	Village string  `json:"village,omitempty"`
}

func (res viewPropRes) Code() int {
	return http.StatusOK
}

func (res viewPropRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewPropRes) Empty() bool {
	return false
}

type propPageRes struct {
	pageRes
	Properties []viewPropRes `json:"properties"`
}

func (res propPageRes) Code() int {
	return http.StatusOK
}

func (res propPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res propPageRes) Empty() bool {
	return false
}

type updateOwnerRes struct {
	ID    string `json:"id,omitempty"`
	Fname string `json:"fname,omitempty"`
	Lname string `json:"lname,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func (res updateOwnerRes) Code() int {
	return http.StatusOK
}

func (res updateOwnerRes) Headers() map[string]string {
	return map[string]string{}
}

func (res updateOwnerRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type viewOwnerRes struct {
	ID    string `json:"id,omitempty"`
	Fname string `json:"fname,omitempty"`
	Lname string `json:"lname,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func (res viewOwnerRes) Code() int {
	return http.StatusOK
}

func (res viewOwnerRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewOwnerRes) Empty() bool {
	return false
}

type ownerPageRes struct {
	pageRes
	Owners []viewOwnerRes `json:"owners"`
}

func (res ownerPageRes) Code() int {
	return http.StatusOK
}

func (res ownerPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res ownerPageRes) Empty() bool {
	return false
}
