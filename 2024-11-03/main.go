package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Hello from %v", request.URL.Path[1:])
	if err != nil {
		log.Fatal("Unable to answer, shutting down")
	}
}

func main() {
	http.HandleFunc("/", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	fmt.Println("Server started")
	log.Fatal(http.ListenAndServe(":8081", nil))
}