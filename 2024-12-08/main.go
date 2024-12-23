package main

import (
	"context"
	"fmt"
	"time"
	// "google.golang.org/grpc"
)

func main() {
	// expect dial time-out on ipv4 blackhole
	_, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	})

	// etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	if err == context.DeadlineExceeded {
		// handle errors
	}

	// etcd clientv3 <= v3.2.9, grpc/grpc-go <= v1.2.1
	// if err == grpc.ErrClientConnTimeout {
	// 	// handle errors
	// }

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()

	spaceCargo := []*Monkey{
		{Name: "George"},
		{Name: "Johnathon The Moon Wizard"},
		{Name: "Crusty Rusty"},
	}

	BunchOfStuff(spaceCargo...)

	// BunchOfStuff(&Monkey{Name: "George"}, &Monkey{Name: "Johnathon The Moon Wizard"})

}


// Functional Options Pattern

type OptFunc func(opts *Option)

type Option struct {
	Port      string
	TLS       bool
	CertFile  string
	KeyFile   string
}

type Server struct {
	opts Option
}

func NewServer(opts ...OptFunc) *Server {
	server := &Server{}
	for _, opt := range opts {
		opt(&server.opts)
	}
	return server
}

type Monkey struct {
	Name string
	YearsInSpace int
}

func (m *Monkey) LaunchIntoSpace() {
	fmt.Println(m.Name, "is in orbit")
	m.YearsInSpace++
}

func BunchOfStuff(monkeys ...*Monkey) {
	for _, monkey := range monkeys {
		monkey.LaunchIntoSpace()
	}
}
