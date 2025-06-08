package main

import (
	"context"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/vmihailenco/msgpack/v5"
)

type Message struct {
	Model string `msgpack:"model"`
}

func main() {
	// Create a Pulsar client
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatalf("Could not create pulsar client: %v", err)
	}
	defer client.Close()

	// Create a producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "model-queue",
	})
	if err != nil {
		log.Fatalf("Could not create producer: %v", err)
	}
	defer producer.Close()

	// Sample models to send
	models := []string{"ModelA", "ModelC", "ModelB", "ModelE", "ModelD", 
		              "ModelG", "ModelF", "ModelI", "ModelH", "ModelJ"}

	ctx := context.Background()

	// Send messages
	for _, model := range models {
		msg := Message{
			Model: model,
		}

		// Encode message using msgpack
		payload, err := msgpack.Marshal(msg)
		if err != nil {
			log.Printf("Error encoding message: %v", err)
			continue
		}

		// Send the message
		msgId, err := producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: payload,
		})
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			continue
		}

		log.Printf("Published message: %v with ID: %v", model, msgId)
		
		// Optional: add small delay between messages
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Finished sending all messages")
}
