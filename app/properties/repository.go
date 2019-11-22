package properties

// Repository defines the api to the properties data store
type Repository interface {
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
