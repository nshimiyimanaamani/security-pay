package ussd

import (
	"context"
	"fmt"

	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/platypus"
)

func (svc *service) Action1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1"

	const menu = "Kwishyura, Andika code y' inzu"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	return platypus.Result{Out: menu, Leaf: end}, nil
}

func (svc *service) Action1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	id, err := params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	property, err := svc.properties.RetrieveByID(ctx, id)
	if err != nil {
		return platypus.Result{Out: "Inzu ntibaho mu nkusanya makuru", Leaf: end}, nil
	}

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	out := fmt.Sprintf(
		"Inzu:%s ya %s %s yishyura:%dRWF\n1. Kwemeza",
		property.ID,
		owner.Fname,
		owner.Lname,
		int(property.Due),
	)

	return platypus.Result{Out: out, Leaf: end}, nil
}

func (svc *service) Action1_1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1"

	const success = "Murakoze"

	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	if err := svc.MakePayment(ctx, property); err != nil {
		return platypus.Result{Out: fail}, nil
	}

	return platypus.Result{Out: success, Leaf: end}, nil
}
