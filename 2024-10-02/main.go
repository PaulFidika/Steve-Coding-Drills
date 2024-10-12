package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	routes := gin.Default()
	routes.GET("/ping", func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Server-Sent Events (SSE) endpoint
	routes.GET("/sse", func(ginCtx *gin.Context) {
		ginCtx.Writer.Header().Set("Content-Type", "text/event-stream")
		ginCtx.Writer.Header().Set("Cache-Control", "no-cache")
		ginCtx.Writer.Header().Set("Connection", "keep-alive")

		ginCtx.Stream(func(w io.Writer) bool {
			for i := 0; i < 10; i++ {
				time.Sleep(1 * time.Second) // Simulate slow response
				fmt.Fprintf(w, "data: Message %d\n\n", i+1)
				ginCtx.Writer.Flush()
			}
			return false
		})
	})

	// Long Polling endpoint
	routes.GET("/longpoll", func(ginCtx *gin.Context) {
		// Simulate waiting for an event
		time.Sleep(1 * time.Second) // Simulate delay for new data

		// Send the response when new data is available
		ginCtx.JSON(http.StatusOK, gin.H{
			"message": "New data available",
		})
	})

	// Chunked Transfer Encoding endpoint
	routes.GET("/chunked", func(ginCtx *gin.Context) {
		ginCtx.Stream(func(w io.Writer) bool {
			for i := 0; i < 10; i++ {
				time.Sleep(1 * time.Second) // Simulate slow response
				fmt.Fprintf(w, "Chunk %d\n", i+1)
				ginCtx.Writer.Flush()
			}
			return false
		})
	})

	routes.Run() // localhost:8080 by default
}
