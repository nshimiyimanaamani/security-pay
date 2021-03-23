package http

import (
	"net/http"
)

// RootResponse ...
type RootResponse struct {
	Title string
}

// APIRoot is the api root handler.
func APIRoot(w http.ResponseWriter, r *http.Request) {}
