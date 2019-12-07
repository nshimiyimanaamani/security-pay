package encoding_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/pkg/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Strings  string
	Integers int
	Floats   float64
	Bytes    []byte
}

func TestEncode(t *testing.T) {
	ts := testStruct{Strings: "test", Integers: 100, Floats: float64(100), Bytes: []byte("test")}

	ctx := context.Background()
	b, err := encoding.Encode(ctx, ts)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	assert.NotNil(t, b, fmt.Sprintf("response bytes should not be nil"))

}

func TestDecode(t *testing.T) {
	ctx := context.Background()

	ts := testStruct{Strings: "test", Integers: 100, Floats: float64(100), Bytes: []byte("test")}
	b, _ := encoding.Encode(ctx, ts)

	var v testStruct
	err := encoding.Decode(ctx, b, &v)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
}
