package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/vmihailenco/msgpack/v5"
)

type Message struct {
	Model string `json:"model"`
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

	jobs := make(chan pulsar.ConsumerMessage, 10)

	// Create a consumer
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "model-queue",
		SubscriptionName: "subscription-name",
		ReceiverQueueSize: 10,
		MessageChannel: jobs,
	})
	if err != nil {
		log.Fatalf("Could not create consumer: %v", err)
	}
	defer consumer.Close()

	// Collect 10 messages
	messages := make([]pulsar.Message, 0, 10)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		msg, err := consumer.Receive(ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}
		messages = append(messages, msg)
		fmt.Printf("Jobs may have gotten this too %+v\n", jobs)
	}

	// Create a slice to hold decoded messages with their original messages
	type MessagePair struct {
		original pulsar.Message
		decoded  Message
	}
	pairs := make([]MessagePair, 0, len(messages))

	// Decode all messages
	for _, msg := range messages {
		var decoded Message
		if err := msgpack.Unmarshal(msg.Payload(), &decoded); err != nil {
			log.Printf("Error decoding message: %v", err)
			consumer.Nack(msg)
			continue
		}
		pairs = append(pairs, MessagePair{
			original: msg,
			decoded:  decoded,
		})
	}

	// Sort messages by Model field
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].decoded.Model < pairs[j].decoded.Model
	})

	// Acknowledge messages in sorted order
	for _, pair := range pairs {
		consumer.Ack(pair.original)
		log.Printf("Acknowledged message with Model: %s", pair.decoded.Model)
	}
}
