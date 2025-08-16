package notificationmanagementservice

import (
	"fmt"
	"net/http"
	"sync"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"
)

type GetNotificationRequest struct {
	UserId   int64
	Page     int64
	PageSize int64
	Search   string
}

type GetNotificationsResponse struct {
	Notifications []schema.Notifications `json:"notifications"`
	Count         int64                  `json:"count"`
}

type ReadNotificationsRequest struct {
	UserId         int64
	NotificationId *int64
}

type ReadNotificationsResponse struct {
	Message string
}

func GetNotification(data GetNotificationRequest) utility.Response[GetNotificationsResponse] {
	offset := (data.Page - 1) * data.PageSize

	var notifications []schema.Notifications
	var count int64
	var notifErr, countErr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		notifErr = db.DB.Model(&schema.Notifications{}).Where("user_id = ?", data.UserId).Offset(int(offset)).Limit(int(data.PageSize)).Find(&notifications).Error
	}()

	go func() {
		defer wg.Done()
		countErr = db.DB.Model(&schema.Notifications{}).Where("user_id = ?", data.UserId).Count(&count).Error
	}()

	wg.Wait()

	if notifErr != nil {
		return utility.Response[GetNotificationsResponse]{
			Success:    false,
			Message:    "failed to fetch notifications",
			Error:      notifErr.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	if countErr != nil {
		return utility.Response[GetNotificationsResponse]{
			Success:    false,
			Message:    "failed to fetch notifications",
			Error:      countErr.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	return utility.Response[GetNotificationsResponse]{
		Success:    true,
		Message:    "Notification fetch successfull",
		StatusCode: http.StatusOK,
		Data: &GetNotificationsResponse{
			Notifications: notifications,
			Count:         count,
		},
	}
}

func ReadNotifications(data ReadNotificationsRequest) utility.Response[ReadNotificationsResponse] {
	query := db.DB.Model(&schema.Notifications{}).Where("user_id = ?", data.UserId)
	if data.NotificationId != nil {
		query.Where("id = ?", data.NotificationId)
	}

	fmt.Println(data.NotificationId)

	err := query.Update("is_read", true).Error

	if err != nil {
		return utility.Response[ReadNotificationsResponse]{
			Success:    false,
			Message:    "failed to read notifications",
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	return utility.Response[ReadNotificationsResponse]{
		Success:    true,
		Message:    "notification read success",
		StatusCode: http.StatusOK,
		Data: &ReadNotificationsResponse{
			Message: "notification read successfully",
		},
	}
}
