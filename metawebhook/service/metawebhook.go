package metawebhookservice

import (
	"net/http"
	"os"
	"visitor-management-system/utility"
)

type VerifyWebhookRequest struct {
	Mode        string
	Challenge   string
	VerifyToken string
}

type VerifyWebhookResponse struct {
	Message string
}

func VerifyWebhook(data VerifyWebhookRequest) utility.Response[VerifyWebhookResponse] {
	if data.Mode == "subscribe" && data.VerifyToken == os.Getenv("META_VERIFY_TOKEN") {
		return utility.Response[VerifyWebhookResponse]{
			Success:    false,
			Message:    "Meta webhook verification failed",
			StatusCode: http.StatusForbidden,
			Data: &VerifyWebhookResponse{
				Message: "Webhook verified",
			},
		}
	} else {
		return utility.Response[VerifyWebhookResponse]{
			Success:    false,
			Message:    "Meta webhook verification failed",
			Error:      "Invalid mode or verify token",
			StatusCode: http.StatusForbidden,
			Data:       nil,
		}
	}
}
