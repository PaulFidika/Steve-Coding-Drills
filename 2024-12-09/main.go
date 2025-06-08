package main

import "fmt"

// Functional Options Pattern, without the Option Struct

func main() {
	server := NewServer(func (s *Server) {
		s.Port = ":9999"
	})

	fmt.Printf("New server: %v\n", server)
	fmt.Printf("New server: %+v\n", server)
}

type OptFunc func(s *Server)

type Server struct {
	Port     string
	TLS      bool
	CertFile string
	KeyFile  string
}

func NewServer(opts ...OptFunc) *Server {
	server := &Server{}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

// func (s *Server) String() string {
// 	return fmt.Sprintf("Server{Port: %s, TLS: %v}", s.Port, s.TLS)
// }