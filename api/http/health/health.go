package health

import "net/http"

// Health indicates the health of the server
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

//Panic simulates panic
func Panic(w http.ResponseWriter, r *http.Request) {
	panic("panic simulation")
}
