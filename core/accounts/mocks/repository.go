package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (accounts.Repository) = (*repositoryMock)(nil)

type repositoryMock struct {
	mu       sync.Mutex
	counter  uint64
	accounts map[string]accounts.Account
}

// NewRepository initiates a mock accounts repository
func NewRepository() accounts.Repository {
	return &repositoryMock{
		accounts: make(map[string]accounts.Account),
	}
}

func (repo *repositoryMock) Save(ctx context.Context, acc accounts.Account) (accounts.Account, error) {
	const op errors.Op = "accounts/repository.Save"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, v := range repo.accounts {
		if v.ID == acc.ID {
			return accounts.Account{}, errors.E(op, "account already exists", errors.KindAlreadyExists)
		}
	}

	repo.counter++
	acc.ID = strconv.FormatUint(repo.counter, 10)

	repo.accounts[acc.ID] = acc

	return acc, nil
}

func (repo *repositoryMock) Update(ctx context.Context, acc accounts.Account) error {
	const op errors.Op = "accounts/repository.Update"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.accounts[acc.ID]; !ok {
		return errors.E(op, "account not found", errors.KindNotFound)
	}

	repo.accounts[acc.ID] = acc

	return nil
}

func (repo *repositoryMock) Retrieve(ctx context.Context, id string) (accounts.Account, error) {
	const op errors.Op = "accounts/repository.Retrieve"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	empty := accounts.Account{}

	account, ok := repo.accounts[id]
	if !ok {
		return empty, errors.E(op, "account not found", errors.KindNotFound)
	}

	return account, nil
}

func (repo *repositoryMock) List(ctx context.Context, offset, limit uint64) (accounts.AccountPage, error) {
	const op errors.Op = "accounts/repository.List"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	items := make([]accounts.Account, 0)

	if offset < 0 || limit <= 0 {
		return accounts.AccountPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	for _, v := range repo.accounts {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := accounts.AccountPage{
		Accounts: items,
		PageMetadata: accounts.PageMetadata{
			Total:  repo.counter,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}
