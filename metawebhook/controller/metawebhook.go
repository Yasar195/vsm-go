package metawebhookcontroller

import (
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
	c.JSON(resp.StatusCode, resp)
}
