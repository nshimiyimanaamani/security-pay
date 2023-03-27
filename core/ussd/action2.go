package ussd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/platypus"
)

func (svc *service) action2(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action2"

	const menu = "Andika nimero yawe ya telephone\n"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}
	return platypus.Result{Out: menu, Leaf: end}, nil

}

func (svc *service) action2_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action2_1"

	const success = "Amazu abanditseho:\n"

	params := platypus.ParamsFromContext(ctx)

	end, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	var owner owners.Owner

	owner.Phone, err = params.GetString("phone")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindNotFound)
	}

	owner, err = svc.owners.RetrieveByPhone(ctx, owner.Phone)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	var page properties.PropertyPage

	page, err = svc.properties.RetrieveByOwner(ctx, owner.ID, 0, 5)
	if err != nil {
		return platypus.Result{Out: "Ntabwo mwanditse mu nkusanya makuru", Leaf: end}, nil
	}

	if len(page.Properties) == 0 {
		return platypus.Result{Out: "Nta mazu abanditseho", Leaf: end}, nil
	}
	return platypus.Result{Out: PrintProperties(success, page), Leaf: end}, nil
}

// PrintProperties ...
func PrintProperties(static string, page properties.PropertyPage) string {
	var buf bytes.Buffer

	_, _ = buf.WriteString(static)

	for n, property := range page.Properties {
		id := property.ID
		sector := property.Address.Sector
		cell := property.Address.Cell
		village := property.Address.Village
		_, _ = buf.WriteString(fmt.Sprintf("%d.%s-%s-%s:'%s'\n", n+1, sector, cell, village, id))
	}
	return buf.String()
}
