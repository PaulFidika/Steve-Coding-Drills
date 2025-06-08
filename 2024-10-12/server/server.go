package main

import (
	"log"
	"net"
	"os"

	chat "steve-fidika/2024-10-12/server/pb/chat"
	pb "steve-fidika/2024-10-12/server/pb/coffeeshop"

	"google.golang.org/grpc"
)

func main() {
	// Step 1: setup coffee server
	listener1, err := net.Listen("tcp", ":9013")
	if err != nil {
		log.Fatalf("Failed to listen to port 9013")
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCoffeeShopServer(grpcServer, &server{})

	// Step 2: setup chat server
	listener2, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen to port 50051: %v", err)
		os.Exit(1)
	}

	grpcServer2 := grpc.NewServer()
	chat.RegisterChatServiceServer(grpcServer2, &chatServer{})

	// Run both servers in parallel
	errChan := make(chan error, 2)

	go func() {
		errChan <- grpcServer.Serve(listener1)
	}()

	go func() {
		errChan <- grpcServer2.Serve(listener2)
	}()

	// Wait for either server to return an error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}