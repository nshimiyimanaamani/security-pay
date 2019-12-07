package redis

import (
	"context"

	"github.com/go-redis/redis/v7"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/encoding"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Queue) = (*queue)(nil)

type queue struct {
	cli *redis.Client
}

// NewQueue initialises the queue
func NewQueue(client *redis.Client) payment.Queue {
	return &queue{client}
}

func (queue *queue) Set(ctx context.Context, tx payment.Transaction) error {
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

func (queue *queue) Get(ctx context.Context, uid string) (payment.Transaction, error) {
	const op errors.Op = "queue.Pop"

	empty := payment.Transaction{}

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

	var tx payment.Transaction
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
