package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error":      "missing auth token",
				"message":    "token validation failed",
				"data":       nil,
				"StatusCode": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error":      "invalid auth token format",
				"message":    "token validation failed",
				"data":       nil,
				"StatusCode": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error":      "invalid or expired token",
				"message":    "token validation failed",
				"data":       nil,
				"StatusCode": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success":    false,
					"error":      "token expired",
					"message":    "token validation failed",
					"data":       nil,
					"StatusCode": http.StatusUnauthorized,
				})
				c.Abort()
				return
			}
			c.Set("claims", claims)
		}

		c.Next()
	}
}
