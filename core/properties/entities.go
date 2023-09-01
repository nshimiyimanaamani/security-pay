package properties

import (
	"time"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// Property defines a property(house) data model
type Property struct {
	ID         string    `json:"id,omitempty"`
	Due        float64   `json:"due,string,omitempty"`
	Owner      Owner     `json:"owner,omitempty"`
	Address    Address   `json:"address,omitempty"`
	Occupied   bool      `json:"occupied,omitempty"`
	ForRent    bool      `json:"for_rent,omitempty"`
	Namespace  string    `json:"namespace"`
	RecordedBy string    `json:"recorded_by,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// Filters defines a filter for property search
type Filters struct {
	Owner     *string
	Names     *string
	Phone     *string
	Sector    *string
	Cell      *string
	Village   *string
	Namespace *string
	Limit     *uint64
	Offset    *uint64
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

// Validate address
func (addr *Address) Validate() error {
	const op errors.Op = "app/properties/address.Validate"

	if addr.Sector == "" || addr.Cell == "" || addr.Village == "" {
		return errors.E(op, "invalid property: invalid address", errors.KindBadRequest)
	}
	return nil
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Amount float64 `json:"amount"`
	Offset uint64
	Limit  uint64
}

// Validate validates a Property entity an returns nil error if it's valid.
func (prt *Property) Validate() error {
	const op errors.Op = "app/properties/property.Validate"
	if err := prt.Owner.Validate(); err != nil {
		return errors.E(op, err, errors.Kind(err))
	}

	if err := prt.Address.Validate(); err != nil {
		return errors.E(op, err, errors.Kind(err))
	}
	if prt.Due == float64(0) {
		return errors.E(op, "invalid property: missing due", errors.KindBadRequest)
	}
	if prt.RecordedBy == "" {
		return errors.E(op, "invalid property: missing recording agent", errors.KindBadRequest)
	}
	if prt.Namespace == "" {
		return errors.E(op, "invalid property: missing namespace tag", errors.KindBadRequest)
	}
	return nil
}

// Owner defines a property owner
type Owner struct {
	ID        string `json:"id,omitempty"`
	Fname     string `json:"fname,omitempty"`
	Lname     string `json:"lname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// OwnerPage ist of owners
type OwnerPage struct {
	Owners []Owner
	PageMetadata
}

// Validate validates owner instance fields
func (ow *Owner) Validate() error {
	const op errors.Op = "app/properties/owner.Validate"
	if ow.ID == "" {
		return errors.E(op, "invalid property: missing owner", errors.KindBadRequest)
	}
	return nil
}
