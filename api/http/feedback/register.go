package feedback

import "github.com/gorilla/mux"

import "net/http"

// RegisterHandlers ...
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(RecodeRoute, Recode(opts.Logger, opts.Service)).Methods(http.MethodPost)
	r.Handle(RetrieveRoute, Retrieve(opts.Logger, opts.Service)).Methods(http.MethodGet)
	r.Handle(DeleteRoute, Delete(opts.Logger, opts.Service)).Methods(http.MethodDelete)
	r.Handle(UpdateRoute, Update(opts.Logger, opts.Service)).Methods(http.MethodPut)
}
