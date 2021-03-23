package encoding

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/vmihailenco/msgpack/v4"
)

// Encode serialises structs into bytes
func Encode(ctx context.Context, v interface{}) ([]byte, error) {
	const op errors.Op = "encoding/binary/Encode"

	b, err := msgpack.Marshal(v)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	return b, nil
}

// Decode deserialises bytes into structs
func Decode(ctx context.Context, b []byte, v interface{}) error {
	const op errors.Op = "encoding/binary/Decode"

	if err := msgpack.Unmarshal(b, v); err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
