package ussd

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/nshimiyimanaamani/paypack-backend/core/owners"
	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/platypus"
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

// For kwishyura ubanze inzu ukoresheje nimero ya telefoni yanditse ariko wishyura ukoresheje indi nimero
func (svc *service) action2_input_2_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.action2_input_2_1"

	const success = "Murakoze gukoresha serivise za PayPack"
	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	status, err := svc.Pay(ctx, property, cmd.Phone)
	if err != nil {
		return platypus.Result{Out: status, Leaf: true}, errors.E(op, err)
	}

	return platypus.Result{Out: success, Leaf: leaf}, nil
}

// Kwemeza Kwishyura ukoresheje option 2 approve to make payment
func (svc *service) action2_input_1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.action2_input_1_1"

	const success = "Murakoze gukoresha serivise za PayPack"
	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	status, err := svc.Pay(ctx, property, property.Owner.Phone)
	if err != nil {
		return platypus.Result{Out: status, Leaf: true}, errors.E(op, err)
	}

	return platypus.Result{Out: success, Leaf: leaf}, nil
}

// Guhitamo nimero yi inzu ushaka kwishurira
func (svc *service) action2_payment_preview(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.action2_payment_preview"

	const fail = "code mwanditse ntibashije kuboneka"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		return platypus.Result{Out: "Nta nimero ya telefoni mwashyizemo", Leaf: true}, errors.E(op, err, errors.KindNotFound)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	var in = "Ugiye kwishyurira inzu ifite '%s' ya %s %s iri mu murenge:'%s' akagari: '%s' umudugudu '%s' yishyura:%dRWF\n1. Kwemeza Kwishyura \n2. Kwemeza Kwishyura amezi menshi"

	invoices, err := svc.invoice.Unpaid(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	if len(invoices.Invoices) > 0 {
		if invoices.Total == 1 {
			in += fmt.Sprintf("\n3. Kwemeza Kwishyura ikirarane cy' ukwezi (%s)"+" angana:%sRWF", invoices.Invoices[0].CreatedAt.Format("2006-01-02"), strconv.Itoa(int(invoices.TotalAmount)))
		} else {
			in += fmt.Sprintf("\n3. Kwemeza Kwishyura ikirarane cy' amezi (%d)"+" angana:%sRWF", invoices.Total, strconv.Itoa(int(invoices.TotalAmount)))
		}
	}

	out := fmt.Sprintf(
		in,
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

// Kureba code za amazu ukoresheje telefoni
func (svc *service) action2_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action2_1"

	const success = "Hitamo kwishyura kumazu yanditse:\n"

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
		switch errors.Kind(err) {
		case errors.KindNotFound:
			return platypus.Result{Out: "Nimero mwashyizemo ntabwo yanditse mu nkusanya makuru", Leaf: true}, errors.E(op, err)
		default:
			return platypus.Result{Out: "Mwihagane habaye ikibazo muri sisiteme", Leaf: true}, errors.E(op, err)
		}
	}

	page, err := svc.properties.RetrieveByOwner(ctx, owner.ID, 0, 10)
	if err != nil {
		return platypus.Result{Out: "Mwihagane habaye ikibazo muri sisiteme", Leaf: end}, errors.E(op, err)
	}

	if len(page.Properties) == 0 {
		return platypus.Result{Out: "Nta mazu abanditseho", Leaf: end}, nil
	}

	return platypus.Result{Out: PrintProperties(success, page), Leaf: end}, nil
}

// For paying ibirarane for option 2
func (svc *service) Action2_phone_house_prev(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action2_phone_house_prev"

	var out = "Ungiye kwishyura ikirarane"
	const fail = "Mwongere mugerageze mukanya habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		return platypus.Result{Out: "Nta nimero ya telefoni mwashyizemo", Leaf: true}, errors.E(op, err, errors.KindNotFound)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	invoices, err := svc.invoice.Unpaid(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	if len(invoices.Invoices) == 0 {
		return platypus.Result{Out: "Ntabirarane mufite bibanditseho", Leaf: true}, nil
	}

	if len(invoices.Invoices) > 0 {
		if invoices.Total == 1 {
			out = fmt.Sprintf("%s cyukwezi kumwe kwa %d bingana na:%d (RWF)\n 1. kwemeza kwishyura", out, invoices.Invoices[0].CreatedAt.Month(), int64(invoices.Invoices[0].Amount))
		} else {
			var amount int64
			for _, invoice := range invoices.Invoices {
				amount += int64(invoice.Amount)
			}
			out = fmt.Sprintf("%s cyamezi %d bigana na:%d (RWF)\n 1. kwemeza kwishyura", out, invoices.Total, amount)
		}
	}

	return platypus.Result{Out: out, Leaf: leaf}, nil
}

// Confirm payment for ibirarane for option 2 ku inimero ya telefoni yanditse kunzu
// For kwishyura ibirarane ukoresheje number yanditse kunzu
func (svc *service) Action2_phone_house_1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action2_phone_house_1_1"

	var out = "Murakoze gukoresha serivisi za paypack"
	const fail = "Mwongere mugerageze mukanya habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		return platypus.Result{Out: "Nta nimero ya telefoni mwashyizemo", Leaf: true}, errors.E(op, err, errors.KindNotFound)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	invoices, err := svc.invoice.Unpaid(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	status, err := svc.CreditPay(ctx, property, property.Owner.Phone, invoices.Invoices)
	if err != nil {
		return platypus.Result{Out: status, Leaf: true}, errors.E(op, err)
	}

	return platypus.Result{Out: out, Leaf: leaf}, nil
}

// Confirm payment for ibirarane for option 2 ku nimero ya telefoni uri gukoresha
// For kwishyura ibirarane ukoresheje number yanditse kunzu
func (svc *service) Action2_phone_house_2_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action2_phone_house_2_1"

	var out = "Murakoze gukoresha serivisi za paypack"
	const fail = "Mwongere mugerageze mukanya habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		return platypus.Result{Out: "Nta nimero ya telefoni mwashyizemo", Leaf: true}, errors.E(op, err, errors.KindNotFound)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	invoices, err := svc.invoice.Unpaid(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	status, err := svc.CreditPay(ctx, property, cmd.Phone, invoices.Invoices)
	if err != nil {
		return platypus.Result{Out: status, Leaf: true}, errors.E(op, err)
	}

	return platypus.Result{Out: out, Leaf: leaf}, nil
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
