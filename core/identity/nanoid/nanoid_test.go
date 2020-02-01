package nanoid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestID(t *testing.T) {
	cfg := &Config{8, "1234567890ABCDEF"}
	provider := New(cfg)

	id := provider.ID()
	require.Equal(t, 8, len(id), fmt.Sprintf("expected id('%s') of length: '%d' got '%d:'", id, 8, len(id)))
}
