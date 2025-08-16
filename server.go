package main

import (
	"context"
	"net/http"
	"os"
	"visitor-management-system/db"
	authenticationcontroller "visitor-management-system/iam/authentication/controller"
	usermanagementcontroller "visitor-management-system/iam/usermanagement/controller"
	middlewares "visitor-management-system/middlewares"
	notificationmanagementcontroller "visitor-management-system/notificationmanagement/controller"
	visitormanagementcontrollers "visitor-management-system/visitormanagement/controller"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	if !isLambda() {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading env files")
		}
	}
	db.ConnectDatabase()
	gin.SetMode(gin.ReleaseMode)

	router := setupRouter()

	ginLambda = ginadapter.NewV2(router)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authenticationcontroller.Login)
		auth.POST("/refresh", middlewares.AuthMiddleware(), authenticationcontroller.TokenRefresh)
	}

	user := router.Group("/api/users")
	user.Use(middlewares.AuthMiddleware())
	{
		user.GET("/", usermanagementcontroller.GetUsers)
		user.POST("/", usermanagementcontroller.CreateUser)
	}

	visitor := router.Group("/api/visitors")
	visitor.Use(middlewares.AuthMiddleware())
	{
		visitor.POST("/", visitormanagementcontrollers.CreateVisitor)
		visitor.GET("/", visitormanagementcontrollers.GetVisitors)
	}

	visits := router.Group("/api/visits")
	visits.Use(middlewares.AuthMiddleware())
	{
		visits.POST("/", visitormanagementcontrollers.CreateVisits)
		visits.GET("/", visitormanagementcontrollers.GetVisits)
	}

	notifications := router.Group("/api/notifications")
	notifications.Use(middlewares.AuthMiddleware())
	{
		notifications.GET("/", notificationmanagementcontroller.GetNotifications)
		notifications.GET("/read", notificationmanagementcontroller.ReadNotifications)
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

	return router
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if isLambda() {
		lambda.Start(Handler)
	} else {
		router := setupRouter()
		router.Run()
	}
}

func isLambda() bool {
	return len(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")) > 0
}
