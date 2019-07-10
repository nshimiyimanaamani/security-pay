package mocks

import (
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/properties"
)

var _ (properties.OwnerStore) = (*ownerStoreMock)(nil)

type ownerStoreMock struct {
	mu      sync.Mutex
	counter uint64
	owners  map[string]properties.Owner
}

// NewOwnerStore instantiates a new OwnerStore mirror.
func NewOwnerStore() properties.OwnerStore {
	return &ownerStoreMock{
		owners: make(map[string]properties.Owner),
	}
}

func (str *ownerStoreMock) Save(owner properties.Owner) (string, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	for _, ow := range str.owners {
		if ow.ID == owner.ID {
			return "", properties.ErrConflict
		}
	}

	str.counter++
	owner.ID = strconv.FormatUint(str.counter, 10)
	str.owners[owner.ID] = owner
	return owner.ID, nil
}

func (str *ownerStoreMock) Update(owner properties.Owner) error {
	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.owners[owner.ID]; !ok {
		return properties.ErrNotFound
	}

	str.owners[owner.ID] = owner
	return nil
}

func (str *ownerStoreMock) Retrieve(id string) (properties.Owner, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	val, ok := str.owners[id]
	if !ok {
		return properties.Owner{}, properties.ErrNotFound
	}

	return val, nil
}

func (str *ownerStoreMock) FindOwner(fname, lname, phone string) (properties.Owner, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	for _, val := range str.owners {
		if val.Fname == fname && val.Lname == lname && val.Phone == phone {
			return val, nil
		}
	}
	return properties.Owner{}, properties.ErrNotFound
}

func (str *ownerStoreMock) RetrieveAll(offset, limit uint64) (properties.OwnerPage, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]properties.Owner, 0)

	if offset < 0 || limit <= 0 {
		return properties.OwnerPage{}, nil
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

	page := properties.OwnerPage{
		Owners: items,
		PageMetadata: properties.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}
