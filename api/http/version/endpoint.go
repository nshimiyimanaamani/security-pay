package version

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/pkg/build"
)

// Build returns larissa build information
func Build(w http.ResponseWriter, r *http.Request) {
	encode(w, build.Data())
}
