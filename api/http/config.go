package http

//Config defines the configuration interface
type Config interface {
	//Fetch returns a configuration entry given it's identifier.
	Fetch(string) interface{}
}

var _ Config = (*config)(nil)

type config struct {
	//Port to bind the web application server to
	Port int

	// ProxyCount isThe number of proxies positioned in front of the API. This is used to interpret
	// X-Forwarded-For headers.
	ProxyCount int
}

func (cfg config) Fetch(entry string) interface{} {
	return nil
}
