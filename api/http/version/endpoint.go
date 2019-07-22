package version

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend"

	"github.com/gorilla/mux"
)

// MakeEndpoint make a version endpoint
func MakeEndpoint(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		response := versionRes{
			Service:   paypack.Service,
			GoVersion: paypack.GoVersion,
			GitCommit: paypack.GitCommit,
			GOOS:      paypack.GOOS,
			GOARCH:    paypack.GOARCH,
		}

		if err := EncodeResponse(w, response); err != nil {
			EncodeError(w, err)
			return
		}
	})
}
