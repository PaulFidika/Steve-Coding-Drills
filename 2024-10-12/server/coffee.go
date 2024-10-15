package main

import (
	"context"
	pb "steve-fidika/2024-10-12/server/pb/coffeeshop"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) StreamOrders(request *pb.StreamRequest, srv grpc.ServerStreamingServer[pb.Order]) error {
	items := []*pb.Item{
		{
			Id: "1",
			Name: "Black Cofee",
		},
		{
			Id: "2",
			Name: "Unicorn Drink",
		},
		{
			Id: "3",
			Name: "Chai Latte",
		},
	}

	for i := range items {
		srv.Send(&pb.Order{
			Items: items[0 : i+1],
		})
	}

	return nil

	// return status.Errorf(codes.Unimplemented, "method GetMenu not implemented")
}

func (s *server) PlaceOrder(context context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
	}, nil

	// return nil, status.Errorf(codes.Unimplemented, "method PlaceOrder not implemented")
}

func (s *server) GetOrderStatus(context context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status: "Waiting forever",
	}, nil

	// return nil, status.Errorf(codes.Unimplemented, "method GetOrderStatus not implemented")
}
