package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestAuditFunc(t *testing.T) {
	exec := postgres.NewExecutor(db)

	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "paypack.developers",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}

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

	saved := saveOwner(t, db, owner)

	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	begin := tools.BeginningOfMonth()

	n := uint64(20)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: saved.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			CreatedAt:  tools.AddMonth(begin, -1),
			UpdatedAt:  tools.AddMonth(begin, -1),
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}
		savePropertyOn(t, db, p)

	}

	var start, end = 0, 10

	exp := 10

	got, err := exec.AuditFunc(context.Background(), start, end)

	require.Nil(t, err, fmt.Sprintf("error %v is not nil", err))

	assert.Equal(t, exp, got, fmt.Sprintf("expected count: %d got %d", exp, got))
}

func lastMonth() time.Time {
	now := time.Now()
	return now.AddDate(0, -1, 0)
}

func TestArchiveFunc(t *testing.T) {
	exec := postgres.NewExecutor(db)

	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "paypack.developers",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}

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

	saved := saveOwner(t, db, owner)

	sector := "Kigomna"
	cell := "Kigeme"
	village := "Tetero"

	begin := tools.BeginningOfMonth()

	n := uint64(20)

	for i := uint64(0); i < n; i++ {
		p := properties.Property{
			ID:    nanoid.New(nil).ID(),
			Owner: properties.Owner{ID: saved.ID},
			Address: properties.Address{
				Sector:  sector,
				Cell:    cell,
				Village: village,
			},
			Namespace:  account.ID,
			CreatedAt:  tools.AddMonth(begin, -1),
			UpdatedAt:  tools.AddMonth(begin, -1),
			Due:        float64(1000),
			RecordedBy: agent.Telephone,
			Occupied:   true,
		}
		savePropertyOn(t, db, p)

	}

	var start, end = 0, 10

	exp := 10

	got, err := exec.ArchiveFunc(context.Background(), start, end)

	require.Nil(t, err, fmt.Sprintf("error %v is not nil", err))

	assert.Equal(t, exp, got, fmt.Sprintf("expected count: %d got %d", exp, got))
}
