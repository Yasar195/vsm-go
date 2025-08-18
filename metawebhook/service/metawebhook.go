package metawebhookservice

import (
	"fmt"
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
	Challenge string
}

func VerifyWebhook(data VerifyWebhookRequest) utility.Response[VerifyWebhookResponse] {
	if data.Mode == "subscribe" && data.VerifyToken == os.Getenv("META_VERIFY_TOKEN") {
		return utility.Response[VerifyWebhookResponse]{
			Success:    true,
			Message:    "Meta webhook verification success",
			StatusCode: http.StatusOK,
			Data: &VerifyWebhookResponse{
				Challenge: data.Challenge,
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

func ReceiveWebhook(data any) utility.Response[any] {
	fmt.Println(data)
	return utility.Response[any]{
		Success:    true,
		Message:    "Message recieved",
		StatusCode: http.StatusOK,
	}
}
