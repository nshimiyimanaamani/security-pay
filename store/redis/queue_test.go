package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueueSet(t *testing.T) {
	queue := redis.NewQueue(redisClient)

	key := uuid.New().ID()

	cases := []struct {
		desc string
		tx   *payment.TxRequest
		err  error
	}{
		{
			desc: "cache unique transacrion",
			tx:   &payment.TxRequest{ID: key},
			err:  nil,
		},
		{
			desc: "cache duplicate transaction",
			tx:   &payment.TxRequest{ID: key},
			err:  nil,
		},
	}
	for _, tc := range cases {
		err := queue.Set(context.Background(), tc.tx)
		assert.Nil(t, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))

	}
}

func TestQueuePop(t *testing.T) {
	queue := redis.NewQueue(redisClient)

	const op errors.Op = "queue.Pop"

	key := uuid.New().ID()

	tsTx := &payment.TxRequest{ID: key}
	err := queue.Set(context.Background(), tsTx)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		key  string
		err  error
	}{
		{
			desc: "get transaction by existing key",
			key:  key,
			err:  nil,
		},
		{
			desc: "get transaction by non existing key",
			key:  uuid.New().ID(),
			err:  errors.E(op, "transaction not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		_, err := queue.Get(context.Background(), tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected '%v' got '%v'\n", tc.desc, tc.err, err))
	}
}

func TestRemove(t *testing.T) {
	queue := redis.NewQueue(redisClient)

	const op errors.Op = "queue.Pop"

	key := uuid.New().ID()

	tsTx := &payment.TxRequest{ID: key}
	err := queue.Set(context.Background(), tsTx)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		key  string
		err  error
	}{
		{
			desc: "delete transaction by existing key",
			key:  key,
			err:  nil,
		},
		{
			desc: "delete transaction by non existing key",
			key:  uuid.New().ID(),
			err:  errors.E(op, "transaction not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		err := queue.Remove(context.Background(), tc.key)
		assert.Nil(t, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}

}
