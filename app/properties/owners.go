package properties

import "github.com/ttacon/libphonenumber"

//Owner defines a property owner
type Owner struct {
	ID       string
	Fname    string
	Lname    string
	Phone    string
	Password string
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

//OwnerStore defines owner storage interface.
type OwnerStore interface {
	// Save adds a new owner to the owner store.
	Save(Owner) (string, error)

	// Updates a given owner entity
	Update(Owner) error

	// Retrieve retrieves an owner given their id
	Retrieve(string) (Owner, error)

	// FindOwner retrieves an owner given their fname, lname and phone.
	FindOwner(string, string, string) (Owner, error)

	// RetrieveAll retrieves a subst of owners.
	RetrieveAll(uint64, uint64) (OwnerPage, error)

	RetrieveByPhone(string) (Owner, error)
}
