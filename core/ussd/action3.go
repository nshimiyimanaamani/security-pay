package ussd

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/platypus"
)

func (svc *service) action3(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action3"

	const menu = "1. Umutekano\n2. Isuku\n"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}
	return platypus.Result{Out: menu, Leaf: end}, nil
}
