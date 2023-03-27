package ussd

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/platypus"
)

func (svc *service) Action0(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action0"

	const menu = "Murakaza neza kuri paypack\nUbufasha: 0785868145\n1. kwishyura\n2. Reba code y' inzu yawe\n3. Gutanga ikibazo\n"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	return platypus.Result{Out: menu, Leaf: end}, nil
}
