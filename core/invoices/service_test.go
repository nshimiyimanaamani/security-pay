package invoices_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nshimiyimanaamani/paypack-backend/core/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/core/invoices/mocks"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var creation = time.Now()

var property = "1"

var invs = map[string]invoices.Invoice{
	property: {ID: 1, Amount: 1000, CreatedAt: creation, UpdatedAt: creation},
}

func newService() invoices.Service {
	repo := mocks.NewRepository(invs)
	opts := &invoices.Options{Repo: repo}
	return invoices.New(opts)
}

func TestRetrieveAll(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/invoices/service.RetrieveAll"

	cases := []struct {
		desc     string
		property string
		months   uint
		size     uint
		err      error
	}{
		{
			desc:     "retrieve invoices for existing property",
			property: property,
			months:   1,
			size:     1,
			err:      nil,
		},
		{
			desc:     "retrieve invoices for non existing property",
			property: "invalid property",
			months:   1,
			size:     0,
			err:      errors.E(op, "property doesn't exists"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.RetrieveAll(ctx, tc.property, tc.months)
		size := uint(len(page.Invoices))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
	}
}

func TestRetrievePending(t *testing.T) {}

func TestRetrievePayed(t *testing.T) {}
