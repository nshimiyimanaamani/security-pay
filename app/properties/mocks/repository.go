package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (properties.Repository) = (*repository)(nil)

type repository struct {
	mu         sync.Mutex
	counter    uint64
	owners     map[string]properties.Owner
	properties map[string]properties.Property
}

// NewRepository creates Repositorymirror
func NewRepository(owners map[string]properties.Owner) properties.Repository {
	return &repository{
		owners:     owners,
		properties: make(map[string]properties.Property),
	}
}

func (str *repository) Save(ctx context.Context, property properties.Property) (properties.Property, error) {
	const op errors.Op = "app/properties/mocks/repository.Save"
	str.mu.Lock()
	defer str.mu.Unlock()

	empty := properties.Property{}

	if _, ok := str.owners[property.Owner.ID]; !ok {
		return empty, errors.E(op, "owner not found", errors.KindNotFound)
	}

	for _, prt := range str.properties {
		if prt.ID == property.ID {
			return empty, errors.E(op, "property already exists", errors.KindAlreadyExists)
		}
	}

	str.counter++
	property.ID = strconv.FormatUint(str.counter, 10)
	str.properties[property.ID] = property
	return property, nil
}

func (str *repository) UpdateProperty(ctx context.Context, property properties.Property) error {
	const op errors.Op = "app/properties/mocks/repository.UpdateProperty"

	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.properties[property.ID]; !ok {
		return errors.E(op, "property not found", errors.KindNotFound)
	}

	str.properties[property.ID] = property

	return nil
}

func (str *repository) RetrieveByID(ctx context.Context, id string) (properties.Property, error) {
	const op errors.Op = "app/properties/mocks/repository.RetrieveByID"

	str.mu.Lock()
	defer str.mu.Unlock()

	val, ok := str.properties[id]
	if !ok {
		return properties.Property{}, errors.E(op, "property not found", errors.KindNotFound)
	}

	return val, nil
}

func (str *repository) RetrieveByOwner(ctx context.Context, owner string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "app/properties/mocks/repository.RetrieveByOwner"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]properties.Property, 0)

	if offset < 0 || limit <= 0 {
		return properties.PropertyPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the property belongs to a given owner
	for _, v := range str.properties {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Owner.ID == owner && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (str *repository) RetrieveBySector(ctx context.Context, sector string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "app/properties/mocks/repository.RetrieveBySector"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]properties.Property, 0)

	if offset < 0 || limit <= 0 {
		return properties.PropertyPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the property belongs to a given owner
	for _, v := range str.properties {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Address.Sector == sector && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (str *repository) RetrieveByCell(ctx context.Context, cell string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "app/properties/mocks/repository.RetrieveByCell"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]properties.Property, 0)

	if offset < 0 || limit <= 0 {
		return properties.PropertyPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the property belongs to a given owner
	for _, v := range str.properties {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Address.Cell == cell && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (str *repository) RetrieveByVillage(ctx context.Context, village string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "app/properties/mocks/repository.RetrieveByVillage"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]properties.Property, 0)

	if offset < 0 || limit <= 0 {
		return properties.PropertyPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the property belongs to a given owner
	for _, v := range str.properties {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Address.Village == village && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (str *repository) RetrieveByRecorder(ctx context.Context, user string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "app/properties/mocks/repository.RetrieveByVillage"

	str.mu.Lock()
	defer str.mu.Unlock()

	items := make([]properties.Property, 0)

	if offset < 0 || limit <= 0 {
		return properties.PropertyPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	//check whether the property belongs to a given owner
	for _, v := range str.properties {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.RecordedBy == user && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  str.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}
