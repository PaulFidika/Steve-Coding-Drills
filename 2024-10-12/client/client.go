package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"steve-fidika/2024-10-12/server/pb/chat"
	pb "steve-fidika/2024-10-12/server/pb/coffeeshop"
	"sync"
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

	// Wait for client2 to finish
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := client2(); err != nil {
			log.Printf("client2 error: %v", err)
		}
	}()

	wg.Wait() // Wait for client2 to finish
}

func client2() error {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect ")
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.SayHello(ctx)
	if err != nil {
		return fmt.Errorf("error creating stream: %w", err)
	}
	defer stream.CloseSend()

	messages := []string{"Hello", "How are you?", "Goodbye"}

	done := make(chan bool)

	go func() error {
		for _, msg := range messages {
			err = stream.Send(&chat.Message{Body: msg})
			if err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}
			log.Printf("Sent: %s", msg)
		}

		time.Sleep(400 * time.Millisecond)

		done <- true

		return nil
	}()

	go func() error {
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					done <- true
					break
				}
				return fmt.Errorf("error receiving response: %w", err)
			}
			log.Printf("Received: %s", resp.GetBody())
		}

		return nil
	}()

	fmt.Println("Sending and receiving messages...")

	<- done

	fmt.Println("Done with all messages")

	return nil
}