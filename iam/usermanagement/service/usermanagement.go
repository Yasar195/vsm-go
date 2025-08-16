package usermanagementservice

import (
	"fmt"
	"net/http"
	"os"
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

type CreateUserInput struct {
	Username  string          `json:"userName" validate:"required"`
	UserEmail string          `json:"userEmail" validate:"required,email"`
	Password  string          `json:"password"`
	UserType  schema.UserType `json:"userType" validate:"required,oneof=staff host"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
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
			Where("user_name ILIKE ?", "%"+data.Search+"%").
			Offset(int(offset)).
			Limit(int(data.PageSize)).
			Find(&users).Error
	}()

	go func() {
		defer wg.Done()
		countErr = db.DB.Model(&schema.Users{}).
			Where("id != ?", data.UserId).
			Where("user_name ILIKE ?", "%"+data.Search+"%").
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

func CreateUser(data CreateUserInput) utility.Response[CreateUserResponse] {

	var user = schema.Users{
		Username:  data.Username,
		UserEmail: data.UserEmail,
		UserType:  data.UserType,
	}

	if data.UserType == "staff" {
		if data.Password == "" {
			return utility.Response[CreateUserResponse]{
				Success:    false,
				Message:    "failed to create user",
				Error:      "staff user requires password",
				Data:       nil,
				StatusCode: http.StatusBadRequest,
			}
		}

		hash, err := utility.HashPassword(data.Password)

		if err != nil {
			return utility.Response[CreateUserResponse]{
				Success:    false,
				Message:    "failed to create user",
				Error:      "passoword hash faile",
				Data:       nil,
				StatusCode: http.StatusInternalServerError,
			}
		}

		user.Password = hash

		emailConfig := utility.EmailConfig{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     587,
			SMTPUsername: os.Getenv("ADMIN_EMAIL"),
			SMTPPassword: os.Getenv("ADMIN_PASSWORD"),
			FromEmail:    os.Getenv("ADMIN_EMAIL"),
		}

		emailService := utility.NewEmailService(emailConfig)

		emailerr := emailService.SendEmail(os.Getenv("ADMIN_EMAIL"), "admin created", fmt.Sprintf("Hi\nNew admin create\n\nemail: %s\npassword: %s", user.UserEmail, user.Password))
		if emailerr != nil {
			return utility.Response[CreateUserResponse]{
				Success:    false,
				Message:    "failed to create user",
				Error:      emailerr.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			}
		}

	}

	if err := db.DB.Create(&user).Error; err != nil {
		return utility.Response[CreateUserResponse]{
			Success:    false,
			Message:    "failed to create user",
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}
	}

	return utility.Response[CreateUserResponse]{
		Success:    true,
		Message:    "user created successfully",
		StatusCode: http.StatusCreated,
		Data: &CreateUserResponse{
			Message: "User created successfully",
		},
	}
}
