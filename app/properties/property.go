package properties

// Property defines a property(house) data model
type Property struct {
	Address
	ID    string
	Owner string
	Due   float64 `json:"due,string,omitempty"`
}

// PropertyPage represents a list of transaction.
type PropertyPage struct {
	PageMetadata
	Properties []Property
}

// Address defines a property location
type Address struct {
	Sector  string
	Cell    string
	Village string
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Validate validates a Property entity an returns nil error if it's valid.
func (prt *Property) Validate() error {
	if prt.Owner == "" || prt.Sector == "" || prt.Cell == "" || prt.Village == "" {
		return ErrInvalidEntity
	}
	if prt.Due == float64(0) {
		return ErrInvalidEntity
	}
	return nil
}

// PropertyStore defines the api to the properties data store
type PropertyStore interface {
	// Save adds a new transactiob to the data store returns nil
	// if the operation is successful or otherwise an error.
	Save(Property) (string, error)

	// UpdateProperty udpates the given property entity's mutable fields.
	UpdateProperty(Property) error

	// RetrieveByID retrieves a property entity  given it's unique id.
	RetrieveByID(string) (Property, error)

	// RetrieveByOwner retrieves the subset of properties that where made using the given owner.
	RetrieveByOwner(string, uint64, uint64) (PropertyPage, error)

	// RetrieveBySector retrieves the subset of properties within a given Sector.
	RetrieveBySector(string, uint64, uint64) (PropertyPage, error)

	// RetrieveByCell retrieves the subset of properteis within a given Cell.
	RetrieveByCell(string, uint64, uint64) (PropertyPage, error)

	// RetrieveByVillage retrieves the subset of properties within a given Village.
	RetrieveByVillage(string, uint64, uint64) (PropertyPage, error)
}
