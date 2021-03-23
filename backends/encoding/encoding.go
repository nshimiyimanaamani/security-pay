package encoding

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Serialize encodes requests objects to json
func Serialize(req interface{}) ([]byte, error) {
	const op errors.Op = "backend.Serialize"

	b, err := json.Marshal(req)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	return b, nil
}

// Deserialize decodes response bodies to appropriate objects
func Deserialize(res io.Reader, v interface{}) error {
	const op errors.Op = "backend.Deserialize"

	data, err := ioutil.ReadAll(res)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
