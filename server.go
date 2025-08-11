package main

import (
	"net/http"
	"visitor-management-system/db"
	authenticationcontroller "visitor-management-system/iam/authentication/controller"
	usermanagementcontroller "visitor-management-system/iam/usermanagement/controller"
	middlewares "visitor-management-system/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading env files")
	}

	db.ConnectDatabase()

	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/login", authenticationcontroller.Login)
		auth.POST("/refresh", authenticationcontroller.TokenRefresh)
	}

	user := router.Group("/users")
	user.Use(middlewares.AuthMiddleware())
	{
		user.GET("/", usermanagementcontroller.GetUsers)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"message":    "No matching resource",
			"error":      "No matching resource found in the server",
			"data":       nil,
			"statusCode": 404,
		})
	})

	router.Run()

}
