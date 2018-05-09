package http

import (
	"log"
	"net"
)

//Server 的封装
type Server struct {
	Response Response
	Request  Request
}

type Config struct {
}

func (srv *Server) Config(config Config) {

}

// Start the server
func (srv *Server) Start() {

	listener, err := net.Listen("tcp", "127.0.0.1:8089")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go srv.Request.Parse(conn)
	}

}
