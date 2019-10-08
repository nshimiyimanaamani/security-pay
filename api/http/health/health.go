package health

import "net/http"

// Health indicates the health of the server
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
