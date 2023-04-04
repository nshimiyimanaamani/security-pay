package ussd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/platypus"
)

// Action01 corresponds to Action1
func (svc *service) Action1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1"

	var menu = "Kwishyura, Hitamo cg wandike code y' inzu:\n"

	params := platypus.ParamsFromContext(ctx)

	cmd.Phone = strings.TrimPrefix(cmd.Phone, "25")
	cmd.Phone = strings.TrimPrefix(cmd.Phone, "+25")

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	owner, err := svc.owners.RetrieveByPhone(ctx, cmd.Phone)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			_, err = svc.agents.RetrieveAgent(ctx, cmd.Phone)
			if err != nil {
				return platypus.Result{Out: "Ntabwo wemerewe gukora iki gikorwa", Leaf: true}, errors.E(op, fmt.Errorf("error: %v:%v", err, cmd.Phone))
			}

			return platypus.Result{Out: "Kwishyura, Andika code y' inzu", Leaf: leaf}, nil

		default:
			return platypus.Result{Out: "Ntabwo mwanditse mu nkusanya makuru", Leaf: true}, errors.E(op, err)
		}
	}

	var page properties.PropertyPage

	page, err = svc.properties.RetrieveByOwner(ctx, owner.ID, 0, 10, owner.Fname)
	if err != nil {
		return platypus.Result{}, nil
	}

	if len(page.Properties) == 0 {
		return platypus.Result{Out: "Nta mazu abanditseho", Leaf: true}, nil
	}

	return platypus.Result{Out: PrintProperties(menu, page), Leaf: leaf}, nil
}

// Choose the payment phone number to be used screen
func (svc *service) Action1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1"

	const success = "Numero yo kwishyura\n1. Iyanditsweho inzu\n2. Iyo uri gukoresha\n"

	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}
	return platypus.Result{Out: success, Leaf: leaf}, nil
}

// ActionPreview corresponds to both ) Action1_1_1_1 and ) Action1_1_1_2
func (svc *service) ActionPreview(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_2_1"

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

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, cmd.Phone)
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

// Option for kwishyura ukoresheje number yanditse kunzu wishyura amaze menshyi
func (svc *service) Action1_1_1_1_2(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1_2"

	const out = "Injiza umubare w' amezi ushaka kwishyura"
	const fail = "Mwongere mugerageze mukanya habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	return platypus.Result{Out: out, Leaf: leaf}, nil
}

// Option for kwishyura ukoresheje number yanditse kunzu wishyura amaze menshyi
func (svc *service) Action1_1_1_1_2_Input(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1_2_Input"

	const out = "Ungiye kwishyura amezi %d angana:%d RWF \n1. Kwemeza Kwishyura"
	const fail = "Mwongere mugerageze mukanya habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	leaf, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	monthVal, err := params.GetInt("input")
	if err != nil {
		return platypus.Result{Out: "Mwashyizemo ibitaribyo hemewe imibare gusa", Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		phone = cmd.Phone
	}

	if phone == " " {
		phone = cmd.Phone
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

	return platypus.Result{Out: fmt.Sprintf(out, monthVal, int(property.Due)*monthVal), Leaf: leaf}, nil
}

// Confirm payment for paying more than one month and include the current month if not paid using number registered on house
func (svc *service) Action1_1_1_1_2_Input_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1_2_Input_1"

	const success = "Murakoze gukoresha serivise za PayPack Mwemeje kwishyura amezi %d angana:%d RWF"
	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	_, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	monthVal, err := params.GetInt("input")
	if err != nil {
		return platypus.Result{Out: "Mwashyizemo ibitaribyo hemewe imibare gusa", Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		phone = cmd.Phone
	}

	if phone == "" {
		phone = cmd.Phone
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

	status, err := svc.BulkPay(ctx, property, property.Owner.Phone, monthVal)
	if err != nil {
		return platypus.Result{Out: status, Leaf: true}, errors.E(op, err)
	}

	return platypus.Result{Out: fmt.Sprintf(success, monthVal, monthVal*int(property.Due)), Leaf: true}, nil
}

// Confirm payment for paying more than one month and include the current month if not paid using current number
func (svc *service) Action1_1_1_2_2_Input_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_2_2_Input_1"

	const success = "Murakoze gukoresha serivise za PayPack Mwemeje kwishyura amezi %d angana:%d RWF"
	const fail = "Mwongere mugerageze habaye ikibazo"

	params := platypus.ParamsFromContext(ctx)

	_, err := params.GetBool("isleaf")
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	monthVal, err := params.GetInt("input")
	if err != nil {
		return platypus.Result{Out: "Mwashyizemo ibitaribyo hemewe imibare gusa", Leaf: true}, errors.E(op, err, errors.KindUnexpected)
	}

	var property properties.Property

	property.ID, err = params.GetString("id")
	if err != nil {
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	phone, err := params.GetString("phone")
	if err != nil {
		phone = cmd.Phone
	}

	if phone == "" {
		phone = cmd.Phone
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

	status, err := svc.BulkPay(ctx, property, cmd.Phone, monthVal)
	if err != nil {
		return platypus.Result{Out: status, Leaf: true}, errors.E(op, err)
	}

	return platypus.Result{Out: fmt.Sprintf(success, monthVal, monthVal*int(property.Due)), Leaf: true}, nil
}

// For paying ibirarane ukoresheje number yanditse kunzu
func (svc *service) Action1_1_1_1_3(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1_3"

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

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, cmd.Phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	invoices, err := svc.invoice.Unpaid(ctx, property.ID)
	if err != nil {
		return platypus.Result{}, errors.E(op, err)
	}

	if len(invoices.Invoices) > 0 {
		if invoices.Total == 1 {
			out = fmt.Sprintf("%s cyukwezi kumwe kwa %d bingana na:%f (RWF)\n 1. kwemeza kwishyura", out, invoices.Invoices[0].CreatedAt.Month(), invoices.Invoices[0].Amount)
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

// For kwishyura ibirarane ukoresheje number yanditse kunzu
func (svc *service) Action1_1_1_1_3_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1_3_1"

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

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, cmd.Phone)
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

// For kwishyura ibirarane ukoresheje number uri gukoresha
func (svc *service) Action1_1_1_2_3_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_2_3_1"

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

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, cmd.Phone)
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

// Option for kwishyura ukoresheje number yanditse kunzu
func (svc *service) Action1_1_1_1(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_1"

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

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, cmd.Phone)
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

// For kwishyura ukoresheje indi number
func (svc *service) Action1_1_1_2(ctx context.Context, cmd *platypus.Command) (platypus.Result, error) {
	const op errors.Op = "core/ussd/service.Action1_1_1_2"

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
		return platypus.Result{}, errors.E(op, err, errors.KindUnexpected)
	}

	// check if the entered input is number and check the corresponding property
	property.ID, err = svc.matchProperty(ctx, property.ID, cmd.Phone)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	property, err = svc.properties.RetrieveByID(ctx, property.ID)
	if err != nil {
		return platypus.Result{Out: fail, Leaf: true}, errors.E(op, err)
	}

	status, err := svc.Pay(ctx, property, cmd.Phone)
	if err != nil {
		return platypus.Result{Out: status, Leaf: leaf}, errors.E(op, err)
	}
	return platypus.Result{Out: success, Leaf: leaf}, nil
}

// Check the code corresponding to the entered index number
func (svc *service) matchProperty(ctx context.Context, id, phone string) (string, error) {

	codeIndex, err := strconv.Atoi(id)
	if err == nil && codeIndex <= 10 {
		phone = strings.TrimPrefix(phone, "25")
		phone = strings.TrimPrefix(phone, "+25")

		owner, err := svc.owners.RetrieveByPhone(ctx, phone)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve owner by phone: %w", err)
		}

		page, err := svc.properties.RetrieveByOwner(ctx, owner.ID, 0, 10, owner.Fname)
		if err != nil {
			return "", err
		}

		if codeIndex > len(page.Properties) || len(page.Properties) == 0 {
			return "", fmt.Errorf("mwihangane ntanzu zibaruye kuri iyi numero")
		}

		for index, prop := range page.Properties {
			if (index + 1) == codeIndex {
				id = prop.ID
				break
			}
		}
	}

	return id, nil
}
