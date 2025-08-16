package notificationmanagementcontroller

import (
	"strconv"
	notificationmanagementservice "visitor-management-system/notificationmanagement/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetNotifications(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	search := c.DefaultQuery("search", "")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 10
	}

	resp := notificationmanagementservice.GetNotification(notificationmanagementservice.GetNotificationRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	c.JSON(resp.StatusCode, resp)

}

func ReadNotifications(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	notificationIdStr := c.Query("notificationId")
	var notificationId *int64

	if notificationIdStr != "" {
		id, err := strconv.ParseInt(notificationIdStr, 10, 64)
		if err == nil {
			notificationId = &id
		}
	}

	resp := notificationmanagementservice.ReadNotifications(notificationmanagementservice.ReadNotificationsRequest{
		UserId:         userID,
		NotificationId: notificationId,
	})

	c.JSON(resp.StatusCode, resp)
}
