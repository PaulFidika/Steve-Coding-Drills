package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Start the Python process
	cmd := exec.Command("python3", "script.py")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Create a reader for user input
	reader := bufio.NewReader(os.Stdin)

	for {
		// Get user input
		fmt.Print("Enter request (or 'quit' to exit): ")
		request, _ := reader.ReadString('\n')
		request = strings.TrimSpace(request)

		if request == "quit" {
			break
		}

		// Send request to Python process
		fmt.Fprintln(stdin, request)

		// Read response from Python process
		response, err := bufio.NewReader(stdout).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		fmt.Println("Response:", strings.TrimSpace(response))
	}

	// Close the stdin to signal the Python process to exit
	stdin.Close()
	cmd.Wait()
}
