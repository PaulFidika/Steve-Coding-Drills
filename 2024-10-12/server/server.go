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
	err = grpcServer.Serve(listener1)

	if err != nil {
		log.Fatalf("Failed to start server on port 9013")
		os.Exit(1)
	}

	// Step 2: setup chat server

	listener2, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen to port 50051: %v", err)
		os.Exit(1)
	}

	grpcServer2 := grpc.NewServer()
	chat.RegisterChatServiceServer(grpcServer2, &chatServer{})
	err = grpcServer2.Serve(listener2)

	if err != nil {
		log.Fatalf("Failed to start grpc server: %v", err)
		os.Exit(1)
	}

}