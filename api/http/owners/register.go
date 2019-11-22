package owners

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterHandlers ....
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(RegisterOwnerRoute, Register(opts.Logger, opts.Service)).Methods(http.MethodPost)

	r.Handle(RetrieveByPhoneRoute, RetrieveByPhone(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("phone", "{phone}")

	r.Handle(SearchOwnerRoute, Search(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("fname", "{fname}", "lname", "{lname}", "phone", "{phone}")

	r.Handle(RetrieveOwnerRoute, Retrieve(opts.Logger, opts.Service)).Methods(http.MethodGet)
	r.Handle(UpdateOwnerRoute, Update(opts.Logger, opts.Service)).Methods(http.MethodPut)

	r.Handle(ListOwnersRoute, List(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("offset", "{offset}", "limit", "{limit}")

}
