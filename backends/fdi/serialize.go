package fdi

import (
	"encoding/json"
	"io"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Serialize encodes requests objects to json
func Serialize(req interface{}) ([]byte, error) {
	const op errors.Op = "fdi.Serialize"

	b, err := json.Marshal(req)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	return b, nil
}

// Deserialize decodes response bodies to appropriate objects
func Deserialize(res io.Reader, v interface{}) error {
	const op errors.Op = "fdi.Deserialize"

	//use json Unmarshal instead

	if err := json.NewDecoder(res).Decode(v); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
