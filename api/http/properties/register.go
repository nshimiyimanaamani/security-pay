package properties

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterHandlers ...
func RegisterHandlers(r *mux.Router, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}

	r.Handle(RegisterPRoute, RegisterProperty(opts.Logger, opts.Service)).Methods(http.MethodPost)
	r.Handle(RetrievePRoute, RetrieveProperty(opts.Logger, opts.Service)).Methods(http.MethodGet)
	r.Handle(UpdatePRoute, UpdateProperty(opts.Logger, opts.Service)).Methods(http.MethodPut)

	r.Handle(ListPRoute, ListPropertyByCell(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("cell", "{cell}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, ListPropertyByOwner(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("owner", "{owner}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, ListPropertyBySector(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("sector", "{sector}", "offset", "{offset}", "limit", "{limit}")

	r.Handle(ListPRoute, ListPropertyByVillage(opts.Logger, opts.Service)).Methods(http.MethodGet).
		Queries("village", "{village}", "offset", "{offset}", "limit", "{limit}")

}
