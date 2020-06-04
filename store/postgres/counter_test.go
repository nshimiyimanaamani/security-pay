package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/tools"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountAuditable(t *testing.T) {
	auditable := postgres.NewCounter(db)

	defer CleanDB(t, db)

	account := accounts.Account{ID: "paypack.developers", Name: "remera", NumberOfSeats: 10, Type: accounts.Devs}

	account = saveAccount(t, db, account)

	agent := users.Agent{
		Telephone: random(15),
		FirstName: "first",
		LastName:  "last",
		Password:  "password",
		Cell:      "cell",
		Sector:    "Sector",
		Village:   "village",
		Role:      users.Dev,
		Account:   account.ID,
	}

	agent = saveAgent(t, db, agent)

	owner := properties.Owner{
		ID:    uuid.New().ID(),
		Fname: "rugwiro",
		Lname: "james",
		Phone: "0784677882",
	}
	owner = saveOwner(t, db, owner)

	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	begin := tools.BeginningOfMonth()

	n := uint64(10)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: owner.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
			CreatedAt:  tools.AddMonth(begin, -1),
			UpdatedAt:  tools.AddMonth(begin, -1),
		}
		savePropertyOn(t, db, p)
	}

	expected := int(n)
	got, err := auditable.Count(context.Background())
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	assert.Equal(t, expected, got, fmt.Sprintf("expected '%d' got '%d'", expected, got))
}
