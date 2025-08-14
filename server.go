package main

import (
	"context"
	"net/http"
	"os"
	"visitor-management-system/db"
	authenticationcontroller "visitor-management-system/iam/authentication/controller"
	usermanagementcontroller "visitor-management-system/iam/usermanagement/controller"
	middlewares "visitor-management-system/middlewares"
	visitormanagementcontrollers "visitor-management-system/visitormanagement/controller"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda

func init() {
	if !isLambda() {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading env files")
		}
	}
	db.ConnectDatabase()
	gin.SetMode(gin.ReleaseMode)

	router := setupRouter()

	ginLambda = ginadapter.New(router)
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
	}

	visitor := router.Group("/api/visitors")
	visitor.Use(middlewares.AuthMiddleware())
	{
		visitor.POST("/", visitormanagementcontrollers.CreateVisitor)
	}

	// Handle 404
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
	v1Request := events.APIGatewayProxyRequest{
		HTTPMethod:      req.RequestContext.HTTP.Method,
		Path:            req.RawPath,
		Resource:        req.RouteKey,
		Headers:         req.Headers,
		Body:            req.Body,
		IsBase64Encoded: req.IsBase64Encoded,
		RequestContext: events.APIGatewayProxyRequestContext{
			HTTPMethod: req.RequestContext.HTTP.Method,
			Path:       req.RawPath,
		},
	}

	if req.PathParameters != nil {
		v1Request.PathParameters = req.PathParameters
	}

	if req.QueryStringParameters != nil {
		v1Request.QueryStringParameters = req.QueryStringParameters
	}

	v1Response, err := ginLambda.ProxyWithContext(ctx, v1Request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "Internal server error"}`,
		}, err
	}

	v2Response := events.APIGatewayV2HTTPResponse{
		StatusCode:      v1Response.StatusCode,
		Headers:         v1Response.Headers,
		Body:            v1Response.Body,
		IsBase64Encoded: v1Response.IsBase64Encoded,
	}

	return v2Response, nil
}

func main() {
	// Check if running in Lambda environment
	if isLambda() {
		lambda.Start(Handler)
	} else {
		// For local development
		router := setupRouter()
		router.Run()
	}
}

// Helper function to detect if running in Lambda
func isLambda() bool {
	// Lambda sets this environment variable
	return len(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")) > 0
}
