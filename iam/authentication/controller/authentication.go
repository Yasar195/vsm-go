package authenticationcontroller

import (
	"net/http"

	authenticationservice "visitor-management-system/iam/authentication/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := authenticationservice.Login(body.Email, body.Password)
	c.JSON(resp.StatusCode, resp)
}

func TokenRefresh(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := authenticationservice.RefreshToken(body.RefreshToken, userID)
	c.JSON(resp.StatusCode, resp)
}
