package main

import (
	"context"
	"io"
	"log"
	"os"
	pb "steve-fidika/2024-10-12/server/pb/coffeeshop"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9013", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect ")
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewCoffeeShopClient(conn)
	 
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	orderStream, err := client.StreamOrders(ctx, &pb.StreamRequest{})
	if err != nil {
		log.Fatal("error getting order stream")
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			resp, err := orderStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}

			if err != nil {
				log.Fatalf("can not receive %v", err)
				continue
			}

			log.Printf("Resp received: %v", resp.Items)
			items = append(items, resp.Items...)
		}
	}()

	<- done

	receipt, err := client.PlaceOrder(ctx, &pb.Order{Items: items})
	if err != nil {
		log.Printf("Unable to place order")
	}
	log.Printf("%v", receipt)

	statusCheck, err := client.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Fatalf("Unable to get order status")
	}
	log.Printf("Status: %v", statusCheck.Status)
}

