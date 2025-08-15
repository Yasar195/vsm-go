package usermanagementservice

import (
	"net/http"
	"sync"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"
)

type GetUserRequest struct {
	UserId   int64
	PageSize int64
	Page     int64
	Search   string
}

type GetUserResponse struct {
	Users []schema.Users `json:"users"`
	Count int64          `json:"count"`
}

func GetUsers(data GetUserRequest) utility.Response[GetUserResponse] {
	offset := (data.Page - 1) * data.PageSize

	var users []schema.Users
	var count int64
	var userErr, countErr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		userErr = db.DB.Model(&schema.Users{}).
			Where("id != ?", data.UserId).
			Where("username ILIKE ?", "%"+data.Search+"%").
			Offset(int(offset)).
			Limit(int(data.PageSize)).
			Find(&users).Error
	}()

	go func() {
		defer wg.Done()
		countErr = db.DB.Model(&schema.Users{}).
			Where("id != ?", data.UserId).
			Where("username ILIKE ?", "%"+data.Search+"%").
			Count(&count).Error
	}()

	wg.Wait()

	if userErr != nil {
		return utility.Response[GetUserResponse]{
			Success:    false,
			Message:    "failed to fetch users",
			Error:      userErr.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	if countErr != nil {
		return utility.Response[GetUserResponse]{
			Success:    false,
			Message:    "failed to fetch user count",
			Error:      countErr.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	return utility.Response[GetUserResponse]{
		Success:    true,
		Message:    "users fetched successfully",
		StatusCode: http.StatusOK,
		Data: &GetUserResponse{
			Users: users,
			Count: count,
		},
	}
}
