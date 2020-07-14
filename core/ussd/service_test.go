package ussd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	ownermocks "github.com/rugwirobaker/paypack-backend/core/owners/mocks"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	propmocks "github.com/rugwirobaker/paypack-backend/core/properties/mocks"
	"github.com/rugwirobaker/paypack-backend/core/ussd"
	"github.com/rugwirobaker/paypack-backend/core/ussd/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const prefix = "*662*104#"

func TestProcess(t *testing.T) {
	owner := owners.Owner{
		Fname: "Karori",
		Lname: "Dan",
		Phone: "0785761000",
	}
	owner, owners := newOwnerRepository(owner)

	property := properties.Property{
		Due: 5000,
		Owner: properties.Owner{
			ID: owner.ID,
		},
		Address: properties.Address{
			Sector:  "Remera",
			Cell:    "Kibagagabaga",
			Village: "Ishema",
		},
	}
	property, properties := newPropertyRepository(property)

	props := []string{property.ID}

	svc := newService(properties, owners, props)

	req := &ussd.Request{
		SessionID:   "session",
		GwRef:       "gwref",
		GwTstamp:    "stamp",
		Msisdn:      "msisdn",
		NetworkCode: "net",
		ServiceCode: "service",
		ServiceID:   "serviceid",
		TenantID:    "tenantid",
	}

	cases := []struct {
		desc     string
		input    string
		expected string
		end      int
		err      error
	}{
		{
			desc:     "action0: main menu",
			input:    "*662*104#",
			expected: "Murakaza neza kuri paypack\n1. kwishyura\n2. reba code y' inzu yawe\n",
			end:      1,
		},
		{
			desc:     "action1: kwishyura",
			input:    "*662*104*1#",
			expected: "Kwishyura, Andika code y' inzu",
			end:      1,
		},
		{
			desc:     "action1_1: guhitamo uwishyura",
			input:    fmt.Sprintf("*662*104*1*%s#", property.ID),
			expected: "Numero yo kwishyura\n 1. Iyanditsweho inzu\n2.Iyo uri gukoresha\n",
			end:      1,
		},
		{
			desc:  "action1_1_1[preview]: Iyanditsweho inzu",
			input: fmt.Sprintf("*662*104*1*%s*1#", property.ID),
			expected: fmt.Sprintf(
				"Ugiye kwishyurira inzu ifite '%s' ya %s %s iri mu murenge: '%s' akagari: '%s' umudugudu '%s' yishyura:%dRWF\n1. Kwemeza",
				property.ID,
				owner.Fname,
				owner.Lname,
				property.Address.Sector,
				property.Address.Cell,
				property.Address.Village,
				int(property.Due),
			),
			end: 1,
		},

		{
			desc:  "action1_1_2[preview]: Iyo uri gukoresha",
			input: fmt.Sprintf("*662*104*1*%s*2#", property.ID),
			expected: fmt.Sprintf(
				"Ugiye kwishyurira inzu ifite '%s' ya %s %s iri mu murenge: '%s' akagari: '%s' umudugudu '%s' yishyura:%dRWF\n1. Kwemeza",
				property.ID,
				owner.Fname,
				owner.Lname,
				property.Address.Sector,
				property.Address.Cell,
				property.Address.Village,
				int(property.Due),
			),
			end: 1,
		},
		{
			desc:     "action1_1_1_1: kwemeza kwishura(registered phone)",
			input:    fmt.Sprintf("*662*104*1*%s*1*1#", property.ID),
			expected: "Murakoze gukoresha serivise za PayPack",
			end:      0,
		},
		{
			desc:     "action1_1_1_2: kwemeza kwishura(user-provided phone)",
			input:    fmt.Sprintf("*662*104*1*%s*2*1#", property.ID),
			expected: "Murakoze gukoresha serivise za PayPack",
			end:      0,
		},
		// action2

		{
			desc:     "action2: kureba code y' inzu",
			input:    "*662*104*2#",
			expected: "Andika nimero yawe ya telephone\n",
			end:      1,
		},
		{
			desc:  "action2_1: kureba code y' inzu",
			input: fmt.Sprintf("*662*104*2*%s#", owner.Phone),
			expected: fmt.Sprintf(
				"Amazu abanditseho:\n1.%s-%s-%s:'%s'\n",
				property.Address.Sector,
				property.Address.Cell,
				property.Address.Village,
				property.ID,
			),
			end: 0,
		},
	}

	for _, tc := range cases {
		req.UserInput = tc.input
		ctx := context.Background()
		res, err := svc.Process(ctx, req)
		require.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.expected, res.Text, fmt.Sprintf("desc:'%s'\n expected: '%s'\n got: '%s'", tc.desc, tc.expected, res.Text))
		assert.Equal(t, tc.end, res.End, fmt.Sprintf("desc:'%s' expected '%d' got '%d'", tc.desc, tc.end, res.End))
	}
}

func ownerFixtures(owner owners.Owner) map[string]properties.Owner {
	owners := make(map[string]properties.Owner)
	owners[owner.ID] = convertOwner(owner)
	return owners
}

func newService(
	ps properties.Repository,
	ows owners.Repository,
	properties []string,
) ussd.Service {
	idp := mocks.NewIdentityProvider()
	opts := &ussd.Options{
		IDP:        idp,
		Prefix:     prefix,
		Properties: ps, Owners: ows,
		Payment: newPaymentService(),
	}
	return ussd.New(opts)
}

func newOwnerRepository(owner owners.Owner) (owners.Owner, owners.Repository) {
	repo := ownermocks.NewRepository()
	owner, _ = repo.Save(context.Background(), owner)
	return owner, repo
}

func newPropertyRepository(property properties.Property) (properties.Property, properties.Repository) {
	repo := propmocks.NewRepository(property.Owner.ID)
	property, _ = repo.Save(context.Background(), property)
	return property, repo
}

func newPaymentService() payment.Service {
	return mocks.NewPaymentService()
}

func identityProvider() identity.Provider {
	return nanoid.New(&nanoid.Config{
		Alphabet: nanoid.Alphabet,
		Length:   nanoid.Length},
	)
}

// must go owner must be defined in core.Owner
func convertOwner(owner owners.Owner) properties.Owner {
	return properties.Owner{
		ID:    owner.ID,
		Fname: owner.Fname,
		Lname: owner.Lname,
		Phone: owner.Phone,
	}
}
