package logmanagementservice

import (
	"fmt"
	"net/http"
	"sync"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"
)

type GetApplicationLogsRequest struct {
	Search   string `json:"search"`
	Page     int64  `json:"page"`
	PageSize int64  `json:"pageSize"`
}

type GetApplicationLogsResponse struct {
	Logs  []schema.Logs `json:"logs"`
	Count int64         `json:"count"`
}

func GetApplicationLogs(data GetApplicationLogsRequest) utility.Response[GetApplicationLogsResponse] {
	var logs []schema.Logs
	var count int64
	var logErr, countErr error
	var wg sync.WaitGroup

	wg.Add(2)

	offset := (data.Page - 1) * data.PageSize

	fmt.Println(offset, data.Page, data.PageSize)

	go func() {
		defer wg.Done()
		logErr = db.DB.Model(&schema.Logs{}).
			Order("created_at DESC").
			Offset(int(offset)).
			Limit(int(data.PageSize)).
			Find(&logs).Error
	}()

	go func() {
		defer wg.Done()
		countErr = db.DB.Model(&schema.Logs{}).
			Count(&count).Error
	}()

	wg.Wait()

	if logErr != nil {
		return utility.Response[GetApplicationLogsResponse]{
			Success:    false,
			Message:    "Failed to fetch logs",
			Data:       nil,
			Error:      logErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if countErr != nil {
		return utility.Response[GetApplicationLogsResponse]{
			Success:    false,
			Message:    "Failed to count logs",
			Data:       nil,
			Error:      countErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[GetApplicationLogsResponse]{
		Success: true,
		Message: "Logs retrieved successfully",
		Data: &GetApplicationLogsResponse{
			Logs:  logs,
			Count: count,
		},
		StatusCode: 200,
	}
}
