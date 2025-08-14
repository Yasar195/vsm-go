package usermanagementservice

import (
	"net/http"
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

	err := db.DB.Model(&schema.Users{}).Where("id != ?", data.UserId).Where("username ILIKE ?", "%"+data.Search+"%").Offset(int(offset)).Limit(int(data.PageSize)).Find(&users).Error

	if err != nil {
		return utility.Response[GetUserResponse]{
			Success:    false,
			Message:    "failed to fetch users",
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	cerr := db.DB.Model(&schema.Users{}).Where("id != ?", data.UserId).Where("username ILIKE ?", "%"+data.Search+"%").Count(&count).Error

	if cerr != nil {
		return utility.Response[GetUserResponse]{
			Success:    false,
			Message:    "failed to fetch users",
			Error:      err.Error(),
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
