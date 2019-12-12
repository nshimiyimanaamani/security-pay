package properties

// Property defines a property(house) data model
type Property struct {
	ID         string  `json:"id,omitempty"`
	Due        float64 `json:"due,string,omitempty"`
	Owner      Owner   `json:"owner,omitempty"`
	Address    Address `json:"address,omitempty"`
	Occupied   bool    `json:"occupied,omitempty"`
	RecordedBy string  `json:"recorded_by,omitempty"`
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
	if addr.Sector == "" || addr.Cell == "" || addr.Village == "" {
		return ErrInvalidEntity
	}
	return nil
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Validate validates a Property entity an returns nil error if it's valid.
func (prt *Property) Validate() error {
	if err := prt.Owner.Validate(); err != nil {
		return err
	}

	if err := prt.Address.Validate(); err != nil {
		return err
	}
	if prt.Due == float64(0) {
		return ErrInvalidEntity
	}
	// if prt.RecordedBy == "" {
	// 	return ErrInvalidEntity
	// }
	return nil
}

//Owner defines a property owner
type Owner struct {
	ID    string `json:"id,omitempty"`
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
func (ow *Owner) Validate() error {
	if ow.ID == "" {
		return ErrInvalidEntity
	}
	return nil
}
