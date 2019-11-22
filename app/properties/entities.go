package properties

import "github.com/ttacon/libphonenumber"

// Property defines a property(house) data model
type Property struct {
	ID      string  `json:"id"`
	Due     float64 `json:"due,string,omitempty"`
	Owner   Owner   `json:"owner,omitempty"`
	Address Address `json:"address,omitempty"`
}

// PropertyPage represents a list of transaction.
type PropertyPage struct {
	PageMetadata
	Properties []Property
}

// Address defines a property location
type Address struct {
	Sector  string `json:"sector,omitempty"`
	Cell    string `json:"cell,omitempty"`
	Village string `json:"village,omitempty"`
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Validate validates a Property entity an returns nil error if it's valid.
func (prt *Property) Validate() error {
	if prt.Owner.ID == "" {
		return ErrInvalidEntity
	}

	if prt.Address.Sector == "" || prt.Address.Cell == "" || prt.Address.Village == "" {
		return ErrInvalidEntity
	}
	if prt.Due == float64(0) {
		return ErrInvalidEntity
	}
	return nil
}

//Owner defines a property owner
type Owner struct {
	ID    string `json:"id"`
	Fname string `json:"fname,omitempty"`
	Lname string `json:"lname,omitempty"`
	Phone string `json:"phone,omitempty"`
}

//OwnerPage ist of owners
type OwnerPage struct {
	Owners []Owner
	PageMetadata
}

// Validate validates owner instance fields
func (own *Owner) Validate() error {
	if own.Fname == "" || own.Lname == "" || own.Phone == "" {
		return ErrInvalidEntity
	}

	num, _ := libphonenumber.Parse(own.Phone, "RW")
	if !libphonenumber.IsValidNumberForRegion(num, "RW") {
		return ErrInvalidEntity
	}
	return nil
}
