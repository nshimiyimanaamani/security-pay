package payment

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(ValidatePaymentRoute, Validate(opts.Logger, opts.Service)).Methods(http.MethodPost)
	r.Handle(InitializePaymentRoute, Initialize(opts.Logger, opts.Service)).Methods(http.MethodPost)
}
