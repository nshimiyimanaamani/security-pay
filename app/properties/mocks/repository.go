package mocks

import (
	"sort"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/properties"
)

var _ (properties.Repository) = (*propertyStoreMock)(nil)

type propertyStoreMock struct {
	mu         sync.Mutex
	counter    uint64
	owner      string
	properties map[string]properties.Property
}

// NewRepository creates Repositorymirror
func NewRepository() properties.Repository {
	return &propertyStoreMock{
		properties: make(map[string]properties.Property),
	}
}

func (str *propertyStoreMock) Save(property properties.Property) (string, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	for _, prt := range str.properties {
		if prt.ID == property.ID {
			return "", properties.ErrConflict
		}
	}

	str.counter++
	property.ID = strconv.FormatUint(str.counter, 10)
	str.properties[property.ID] = property
	return property.ID, nil
}

func (str *propertyStoreMock) UpdateProperty(property properties.Property) error {
	str.mu.Lock()
	defer str.mu.Unlock()

	if _, ok := str.properties[property.ID]; !ok {
		return properties.ErrNotFound
	}

	str.properties[property.ID] = property

	return nil
}

func (str *propertyStoreMock) RetrieveByID(id string) (properties.Property, error) {
	str.mu.Lock()
	defer str.mu.Unlock()

	val, ok := str.properties[id]
	if !ok {
		return properties.Property{}, properties.ErrNotFound
	}

	return val, nil
}

func (str *propertyStoreMock) RetrieveByOwner(owner string, offset, limit uint64) (properties.PropertyPage, error) {
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

func (str *propertyStoreMock) RetrieveBySector(sector string, offset, limit uint64) (properties.PropertyPage, error) {
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

func (str *propertyStoreMock) RetrieveByCell(cell string, offset, limit uint64) (properties.PropertyPage, error) {
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

func (str *propertyStoreMock) RetrieveByVillage(village string, offset, limit uint64) (properties.PropertyPage, error) {
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
