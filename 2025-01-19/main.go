package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// Create a rate limiter that allows 10 events per second
	// The second parameter (burst) determines the maximum burst size
	limiter := rate.NewLimiter(1, 10)

	// Simulate a function that processes requests
	processRequest := func(id int) {
		fmt.Printf("Processing request %d at %s\n", id, time.Now().Format(time.RFC3339))
	}

	// Use the limiter to control the rate of requests
	for i := 1; i <= 100; i++ {
		// Introduce a small random delay between requests
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

		// Check if the request is allowed immediately
		if !limiter.Allow() {
			fmt.Printf("Request %d denied due to rate limiting at %s\n", i, time.Now().Format(time.RFC3339))
			continue
		}

		// Process the request
		processRequest(i)
	}
}
