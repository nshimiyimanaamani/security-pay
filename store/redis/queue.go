package redis

import (
	"context"

	"github.com/go-redis/redis/v7"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/encoding"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (payment.Queue) = (*queue)(nil)

type queue struct {
	cli *redis.Client
}

// NewQueue initialises the queue
func NewQueue(client *redis.Client) payment.Queue {
	return &queue{client}
}

func (queue *queue) Set(ctx context.Context, tx *payment.TxRequest) error {
	const op errors.Op = "queue.Set"

	b, err := encoding.Encode(ctx, tx)
	if err != nil {
		return errors.E(op, err)
	}

	_, err = queue.cli.SetNX(tx.ID, b, payment.TxExpiration).Result()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}

func (queue *queue) Get(ctx context.Context, uid string) (payment.TxRequest, error) {
	const op errors.Op = "queue.Pop"

	empty := payment.TxRequest{}

	res, err := queue.cli.Get(uid).Result()
	if err != nil {
		if err == redis.Nil {
			return empty, errors.E(op, "transaction not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	_, err = queue.cli.Del(uid).Result()
	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	var tx payment.TxRequest
	if err := encoding.Decode(ctx, []byte(res), &tx); err != nil {
		return empty, errors.E(op, err, "unable to deserialize transaction", errors.KindUnexpected)
	}
	return tx, nil
}

func (queue *queue) Remove(ctx context.Context, uid string) error {
	const op errors.Op = "queue.Remove"

	_, err := queue.cli.Del(uid).Result()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}
