package main

import (
	"fmt"
	"sync"
	"time"
)

var m = sync.Mutex{}
var wg = sync.WaitGroup{}
var dbData = []string{"id1", "id2", "id3", "id4", "id5", "id6"}
var dbDataMap = map[string]string{
	"id1": "RandomString1",
	"id2": "AnotherRandomString",
	"id3": "YetAnotherString",
	"id4": "FourthRandomString",
	"id5": "FifthRandomString",
	"id6": "LastRandomString",
}

func main() {
	results := make(chan string)
	// done := make(chan bool)
	expectedResults := len(dbData)
	// moreResults := []string{}
	myCounter := 1

	t0 := time.Now()
	for _, id := range dbData {
		wg.Add(1)
		go func(id string) {
			results <- dbCall(id)
			// Without this lock, this block has inconsistent behavior
			// m.Lock()
			var thing = myCounter
			time.Sleep(100_000_000)
			myCounter = thing + 1
			// m.Unlock()
			wg.Done()
		}(id)
	}

	// go func() {
	// 	wg.Wait()
	// 	close(results)
	// }()

	i := 0
	for result := range results {
        fmt.Println("Received result:", result)
		i++
		if i == expectedResults {
			close(results)
		}
    }

	wg.Wait()
    // close(results)

	fmt.Printf("\nCounter results: %v", myCounter)
	fmt.Printf("\nTotal exeuction time: %v", time.Since(t0))
}

func dbCall(id string) string {
	// var delay float32 = rand.Float32() * 2000
	var delay = 2000
	time.Sleep(time.Duration(delay) * time.Millisecond)
	fmt.Println("The result from the database is:", id)
	return dbDataMap[id]
}