package feedback

import (
	"net/http"
	"time"
)

// RecordRes ...
type RecordRes struct {
	ID        string    `json:"id"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty"`
}

// Code returns the https status code
func (res RecordRes) Code() int {
	return http.StatusCreated
}

// Headers returns headers
func (res RecordRes) Headers() map[string]string {
	return map[string]string{}
}

// Empty checks whether the response is empty
func (res RecordRes) Empty() bool {
	return false
}

// RetrieveRes ...
type RetrieveRes struct {
	ID        string    `json:"id"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty"`
}

// Code returns the https status code
func (res RetrieveRes) Code() int {
	return http.StatusOK
}

// Headers returns headers
func (res RetrieveRes) Headers() map[string]string {
	return map[string]string{}
}

// Empty checks whether the response is empty
func (res RetrieveRes) Empty() bool {
	return false
}

// UpdateRes ...
type UpdateRes struct {
	ID        string    `json:"id"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty"`
}

// Code returns the https status code
func (res UpdateRes) Code() int {
	return http.StatusOK
}

// Headers returns headers
func (res UpdateRes) Headers() map[string]string {
	return map[string]string{}
}

// Empty checks whether the response is empty
func (res UpdateRes) Empty() bool {
	return false
}
