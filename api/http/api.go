package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/users"
)

//API defines the application's  http API.
type API struct {
	Users users.Service
}

//statusCodeRecorder extends the functionalities of http.ResponseWriter
type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

//New instatiates a new http API.
func New(users users.Service) *API {
	return &API{
		Users: users,
	}
}

func (api *API) handler(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		if err := checkContentType(r); err != nil {
			encodeError(w, err)
			return
		}

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
		}()

		if err := f(w, r); err != nil {
			encodeError(w, err)
			return
		}

	})
}

//Init initializers the API endpoints.
func (api *API) Init(r *mux.Router) {
	r.Handle("/users", api.handler(api.UserRegisterEndpoint)).Methods("POST")
	r.Handle("/tokens", api.handler(api.UserLoginEndpoint)).Methods("POST")
}

func (api *API) ipAddressFromRequest(r *http.Request) string {
	return ""
}
