package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/owners"
)

var _ (owners.Repository) = (*ownerRepoMock)(nil)

type ownerRepoMock struct {
	mu      sync.Mutex
	counter uint64
	owners  map[string]owners.Owner
}

// NewOwnerRepository instantiates a new Repository mirror.
func NewOwnerRepository() owners.Repository {
	return &ownerRepoMock{
		owners: make(map[string]owners.Owner),
	}
}

func (str *ownerRepoMock) Save(ctx context.Context, owner owners.Owner) (owners.Owner, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	for _, ow := range str.owners {
		if ow.ID == owner.ID {
			return owners.Owner{}, owners.ErrConflict
		}
	}

	str.counter++
	owner.ID = strconv.FormatUint(str.counter, 10)
	str.owners[owner.ID] = owner
	return owner, nil
}

func (str *ownerRepoMock) Update(ctx context.Context, owner owners.Owner) error {
	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.owners[owner.ID]; !ok {
		return owners.ErrNotFound
	}

	str.owners[owner.ID] = owner
	return nil
}

func (str *ownerRepoMock) Retrieve(ctx context.Context, id string) (owners.Owner, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	val, ok := str.owners[id]
	if !ok {
		return owners.Owner{}, owners.ErrNotFound
	}

	return val, nil
}

func (str *ownerRepoMock) Search(ctx context.Context, owner owners.Owner) (owners.Owner, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	for _, val := range str.owners {
		if val.Fname == owner.Fname && val.Lname == owner.Lname && val.Phone == owner.Phone {
			return val, nil
		}
	}
	return owners.Owner{}, owners.ErrNotFound
}

func (str *ownerRepoMock) RetrieveAll(ctx context.Context, offset, limit uint64) (owners.OwnerPage, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]owners.Owner, 0)

	if offset < 0 || limit <= 0 {
		return owners.OwnerPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the tranaction belongs to a given property
	for _, v := range str.owners {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := owners.OwnerPage{
		Owners: items,
		PageMetadata: owners.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (str *ownerRepoMock) RetrieveByPhone(ctx context.Context, phone string) (owners.Owner, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	for _, ow := range str.owners {
		if ow.Phone == phone {
			return ow, nil
		}
	}

	return owners.Owner{}, owners.ErrNotFound
}
