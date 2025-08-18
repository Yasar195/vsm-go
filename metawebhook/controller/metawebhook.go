package metawebhookcontroller

import (
	"net/http"
	metawebhookservice "visitor-management-system/metawebhook/service"

	"github.com/gin-gonic/gin"
)

func VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	challenge := c.Query("hub.challenge")
	veriifyToken := c.Query("hub.verify_token")

	resp := metawebhookservice.VerifyWebhook(metawebhookservice.VerifyWebhookRequest{
		Mode:        mode,
		Challenge:   challenge,
		VerifyToken: veriifyToken,
	})
	c.String(resp.StatusCode, resp.Data.Challenge)
}

func ReceiveWebhook(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	resp := metawebhookservice.ReceiveWebhook(body)
	c.JSON(resp.StatusCode, resp)
}
