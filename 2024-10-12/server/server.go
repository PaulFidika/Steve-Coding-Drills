package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", "50051")
	if err != nil {
		log.Fatalf("Failed to listen to port 50051: %v", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
		os.Exit(1)
	}

}