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

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	return platypus.Result{Out: menu, Leaf: leaf}, nil
}

func (svc *service) Action1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1"

	const success = "Numero yo kwishyura\n 1. Iyanditsweho inzu\n 2. Iyo uri gukoresha\n"

	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}
	return platypus.Result{Out: success, Leaf: leaf}, nil
}

// ActionPreview corresponds to both ) Action1_1_1_1 and ) Action1_1_1_2
func (svc *service) ActionPreview(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_2_1"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}
	out := fmt.Sprintf(
		"Ugiye kwishyurira inzu ifite '%s' ya %s %s iri mu murenge: '%s' akagari: '%s' umudugudu '%s' yishyura:%dRWF\n1. Kwemeza",
		property.ID,
		owner.Fname,
		owner.Lname,
		property.Address.Sector,
		property.Address.Cell,
		property.Address.Village,
		int(property.Due),
	)
	return platypus.Result{Out: out, Leaf: leaf}, nil
}

func (svc *service) Action1_1_1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1"

	const success = "Murakoze gukoresha serivise za PayPack"
	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	status, err := svc.Pay(ctx, property, owner.Phone)
	if err != nil {
		return platypus.Result{Out: status, Leaf: leaf}, nil
	}

	return platypus.Result{Out: status, Leaf: leaf}, nil
}

func (svc *service) Action1_1_1_2(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_2"

	const success = "Murakoze gukoresha serivise za PayPack"
	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	status, err := svc.Pay(ctx, property, cmd.Phone)
	if err != nil {
		return platypus.Result{Out: status, Leaf: leaf}, nil
	}
	return platypus.Result{Out: status, Leaf: leaf}, nil
}
