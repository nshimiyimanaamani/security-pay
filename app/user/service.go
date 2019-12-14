package user

// Service defines users usecases
type Service interface{
	Agents
	Developers
	Managers
}
// Agents ...
type Agents interface{
	RegisterAgent()error
}

// Developers ...
type Developers interface{
	RegisterDeveloper()error
}

// Managers ...
type Managers interface{
	RegisterManager()error
}