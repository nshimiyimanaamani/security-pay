package http

//Server Defines an http Server. I holds transport level state
type Server struct{}

//Options defines optional server configuration.
type Options struct{}

//New create anew http Server instance
/**
 * @todo Add New Http Server Function
 * @body * includes:
		 - http server config definition
		 - options
*/
func New() *Server {
	return &Server{}
}

//Run start the http Server
/**
 * @todo Add Run Http Server Function
 * @body * includes:
		 - define the Http API
*/
func (srv *Server) Run() error { return nil }
