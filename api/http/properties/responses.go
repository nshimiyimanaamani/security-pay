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
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Sector  string `json:"sector"`
	Cell    string `json:"cell"`
	Village string `json:"village"`
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
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Sector  string `json:"sector"`
	Cell    string `json:"cell"`
	Village string `json:"village"`
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
