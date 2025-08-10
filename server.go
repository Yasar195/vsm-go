package main

import (
	"visitor-management-system/db"
	authenticationcontroller "visitor-management-system/iam/authentication/controller"

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
	}

	router.Run()

}
