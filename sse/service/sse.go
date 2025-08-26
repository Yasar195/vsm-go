package sseservice

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetEvents(c *gin.Context) {
	client := c.Request.Context().Done()

	for {
		select {
		case <-client:
			fmt.Println("Client disconnected")
			return
		default:
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(c.Writer, "data: Current time: %s\n\n", timestamp)

			c.Writer.Flush()

			time.Sleep(1 * time.Second)
		}
	}
}
