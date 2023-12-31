package owners

import (
	"errors"
	"regexp"
)

// Sentinel Errors
var (
	ErrInvalidEntity = errors.New("invalid owner entity")
	ErrNotFound      = errors.New("owner entity not found")
	ErrConflict      = errors.New("owner already exists")
)

// Owner defines a property owner
type Owner struct {
	ID    string `json:"id"`
	Fname string `json:"fname,omitempty"`
	Lname string `json:"lname,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// OwnerPage ist of owners
type OwnerPage struct {
	Owners       []Owner `json:"owners"`
	PageMetadata `json:"meta"`
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Validate validates owner instance fields
func (own *Owner) Validate() error {
	if own.Fname == "" || own.Lname == "" || own.Phone == "" {
		return ErrInvalidEntity
	}

	var phoneRegex = `^(\+?25)?(078|079|073|072)\d{7}$`
	if !regexp.MustCompile(phoneRegex).MatchString(own.Phone) {
		return errors.New("invalid phone number provided")
	}

	return nil
}
