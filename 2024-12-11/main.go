package main

import (
	"fmt"

	"github.com/surrealdb/surrealdb.go"
)

type Person struct {
	Name 	string
	Age 	int
}

// https://surrealdb.com/docs/sdk/golang/methods

func main() {
    // Create a connection
    db, err := surrealdb.New("ws://localhost:8000")
    if err != nil {
        panic(err)
    }

    // Start a live query
    queryId, err := surrealdb.Live(db, "SELECT * FROM person", true)
    if err != nil {
        panic(err)
    }
	defer surrealdb.Kill(db, queryId.String())

	channel, _ := db.LiveNotifications(queryId.String())
	
	for {
		val, ok := <-channel
		if !ok {
			// Channel is closed
			return
		}
		// val.ID.Format()
		fmt.Printf("Received: %v\n", val)
	}
}