package transactions

import (
	"net/http"
	"time"

	transport "github.com/rugwirobaker/paypack-backend/api/http"
)

var _ transport.Response = (*recordTransRes)(nil)

type recordTransRes struct {
	ID string `json:"id,omitempty"`
}

func (res recordTransRes) Code() int {
	return http.StatusCreated
}

func (res recordTransRes) Headers() map[string]string {
	return map[string]string{}
}

func (res recordTransRes) Empty() bool {
	return false
}

type viewTransRes struct {
	ID           string            `json:"id,omitempty"`
	Property     string            `json:"property,omitempty"`
	Owner        string            `json:"owner,omitempty"`
	Amount       string            `json:"amount,omitempty"`
	Address      map[string]string `json:"address,omitempty"`
	Method       string            `json:"method,omitempty"`
	DateRecorded time.Time         `json:"recorded,omitempty"`
}

func (res viewTransRes) Code() int {
	return http.StatusOK
}

func (res viewTransRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewTransRes) Empty() bool {
	return false
}

type transPageRes struct {
	pageRes
	Transactions []viewTransRes `json:"transactions"`
}

func (res transPageRes) Code() int {
	return http.StatusOK
}

func (res transPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res transPageRes) Empty() bool {
	return false
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}
