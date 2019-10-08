package payment

import (
	"encoding/json"
	"io"
)

func decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
