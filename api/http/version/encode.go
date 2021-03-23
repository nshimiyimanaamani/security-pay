package version

import (
	"encoding/json"
	"net/http"
)

var contentType = "application/json"

func encode(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(v)
}
