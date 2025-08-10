package authenticationservice

import (
	"os"
	"time"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type LoginResponse struct {
	Tokens   Tokens
	UserName string
	Email    string
}

func Login(email string, password string) utility.Response[LoginResponse] {
	var user schema.Users

	if err := db.DB.Where("user_email = ?", email).First(&user).Error; err != nil {
		return utility.Response[LoginResponse]{
			Success:    false,
			Message:    "Invalid email or password",
			Error:      err.Error(),
			StatusCode: 401,
			Data:       nil,
		}
	}

	if !utility.ComparePassword(user.Password, password) {
		return utility.Response[LoginResponse]{
			Success:    false,
			Message:    "Invalid email or password",
			Error:      "password mismatch",
			StatusCode: 401,
			Data:       nil,
		}
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	refreshsecret := os.Getenv("JWT_SECRET")
	if secret == "" {
		refreshsecret = "dev-refreshsecret"
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.UserEmail,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return utility.Response[LoginResponse]{
			Success:    false,
			Message:    "Failed to generate access token",
			Error:      err.Error(),
			StatusCode: 500,
			Data:       nil,
		}
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.UserEmail,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshsecret))

	return utility.Response[LoginResponse]{
		Success: true,
		Data: &LoginResponse{
			Tokens: Tokens{
				AccessToken:  accessTokenString,
				RefreshToken: refreshTokenString,
			},
			UserName: user.Username,
			Email:    user.UserEmail,
		},
		Message:    "Login successful",
		StatusCode: 200,
	}
}
