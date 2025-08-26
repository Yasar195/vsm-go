package ssecontroller

import (
	sseservice "visitor-management-system/sse/service"

	"github.com/gin-gonic/gin"
)

func GetEvents(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	sseservice.GetEvents(c)
}
