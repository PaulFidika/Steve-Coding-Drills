package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigchan := make(chan os.Signal, 3)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<- sigchan
		fmt.Println("Exiting...")
		os.Exit(3)
	}()

	// Block indefinitely
	select {}
}