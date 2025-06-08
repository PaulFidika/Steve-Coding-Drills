package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Set up channel to listen for interrupt signal
	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Channel to receive lines from stdin
	lineChan := make(chan string, 1)

	// Goroutine to handle the interrupt signal
	// go func() {
	// 	<-sigChan
	// 	fmt.Println("\nShutting down...")
	// 	os.Exit(0)
	// }()
	
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "You fucked up, you fucked up")
		os.Exit(1)
	}

	old, new := os.Args[1], os.Args[2]

	// Goroutine to read lines from stdin
	go func() {
		scan := bufio.NewScanner(os.Stdin)
		for scan.Scan() {
			lineChan <- scan.Text()
		}
		close(lineChan)
	}()

	for {
		select {
		case new_text, ok := <-lineChan:
			if !ok {
				// Channel closed, exit the loop
				return
			}
			if new_text == "exit" {
				os.Exit(0)
			}

			s := strings.Split(new_text, old)
			t := strings.Join(s, new)

			fmt.Println(t)

		case sig := <-sigChan:
			fmt.Printf("\nReceived %s, shutting down gracefully...\n", sig)
			os.Exit(0)
		}
	}

	// Block indefinitely
	select {}
}