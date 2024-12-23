package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const url = "https://jsonplaceholder.typicode.com"

type Todo struct {
	UserId 		int 	`json:"userId"`
	Id 			int 	`json:"Id"`
	Title 		string 	`json:"title"`
	Completed 	bool 	`json:"completed"`
}

func main() {
	// resp, err := http.Get("http://localhost:8080/" + os.Args[1])
	resp, err := http.Get(url + "/todos/1")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

		// Add this debug line
	fmt.Printf("Protocol: %s\n", resp.Proto)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		var item Todo

		err = json.Unmarshal(body, &item)

		if err != nil {
			logger := log.Default()
			logger.Print(err)
		}

		fmt.Println(item)

		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&item); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// _ := template.New("mine")
		// tmpl.Parse(form)
	}
}