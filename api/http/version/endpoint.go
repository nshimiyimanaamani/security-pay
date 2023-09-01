package version

import (
	"net/http"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/build"
)

// Build returns paypack build information
func Build(w http.ResponseWriter, r *http.Request) {
	encode(w, build.Data())
}
