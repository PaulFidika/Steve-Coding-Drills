package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2_000*time.Millisecond)

	defer cancel()

	results := make(chan result)
	list := []string {
		"https://google.com",
		"https://amazon.com",
		"https://facebook.com",
		"https://doujins.com",
		"https://cozy.art",
		"https://fakku.net",
		"https://fal.ai",
	}

	for _, url := range(list) {
		go get(ctx, url, results)
	}

	for range 7 {
		result := <-results

		if result.success {
			fmt.Println("it worked", result.url)
		} else {
			fmt.Println("it failed", result.url)
		}
	}
}

type result struct {
	success bool
	res *http.Response
	url string
}

func get(ctx context.Context, url string, ch chan<- result) {
	start := time.Now()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if res, err := http.DefaultClient.Do(req); err != nil {
		elapsed := time.Since(start)
		fmt.Printf("Request to %s failed after %v\n", url, elapsed)
		ch <- result{false, nil, url}
	} else {
		elapsed := time.Since(start)
		fmt.Printf("Request to %s succeeded in %v\n", url, elapsed)
		ch <- result{true, res, url}
		res.Body.Close()
	}
}