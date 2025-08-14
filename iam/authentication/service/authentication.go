package authenticationservice

import (
	"fmt"
	"os"
	"time"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	Tokens   Tokens `json:"tokens"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func Login(email string, password string) utility.Response[LoginResponse] {
	var user schema.Users

	if err := db.DB.Where("user_email = ? AND user_status = ?", email, "active").First(&user).Error; err != nil {
		return utility.Response[LoginResponse]{
			Success:    false,
			Message:    "Login failed",
			Error:      err.Error(),
			StatusCode: 401,
			Data:       nil,
		}
	}

	if !utility.ComparePassword(user.Password, password) {
		return utility.Response[LoginResponse]{
			Success:    false,
			Message:    "Login failed",
			Error:      "password mismatch",
			StatusCode: 401,
			Data:       nil,
		}
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	refreshsecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshsecret == "" {
		refreshsecret = "dev-refreshsecret"
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.UserEmail,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return utility.Response[LoginResponse]{
			Success:    false,
			Message:    "Login failed",
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

func RefreshToken(refreshToken string, userId int64) utility.Response[RefreshTokenResponse] {
	var refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
	parsedToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshSecret, nil
	})

	if err != nil {
		return utility.Response[RefreshTokenResponse]{
			Success:    false,
			Message:    "Token refresh failed",
			Error:      err.Error(),
			Data:       nil,
			StatusCode: 400,
		}
	}

	if !parsedToken.Valid {
		return utility.Response[RefreshTokenResponse]{
			Success:    false,
			Message:    "Token refresh failed",
			Data:       nil,
			Error:      "Invalid refresh token",
			StatusCode: 400,
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return utility.Response[RefreshTokenResponse]{
			Success:    false,
			Message:    "Token refresh failed",
			Data:       nil,
			Error:      "Invalid token claims",
			StatusCode: 400,
		}
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return utility.Response[RefreshTokenResponse]{
				Success:    false,
				Message:    "Token refresh failed",
				Data:       nil,
				Error:      "refresh token expired",
				StatusCode: 401,
			}
		}
	}

	email, ok := claims["email"].(string)
	userID, ok := claims["user_id"]

	if userID == userId {
		return utility.Response[RefreshTokenResponse]{
			Success:    false,
			Message:    "Token refresh failed",
			Error:      "user mismatch",
			StatusCode: 500,
			Data:       nil,
		}
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return utility.Response[RefreshTokenResponse]{
			Success:    false,
			Message:    "Token refresh failed",
			Error:      err.Error(),
			StatusCode: 500,
			Data:       nil,
		}
	}

	return utility.Response[RefreshTokenResponse]{
		Success: true,
		Data: &RefreshTokenResponse{
			AccessToken: accessTokenString,
		},
		Message:    "Token refresh success",
		StatusCode: 201,
	}

}
