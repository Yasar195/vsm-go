package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
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

// Lambda handler function
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Debug logging - enhanced
	fmt.Printf("Original Request - Method: '%s', Path: '%s', Resource: '%s'\n", req.HTTPMethod, req.Path, req.Resource)
	fmt.Printf("RequestContext - Method: '%s', Path: '%s'\n", req.RequestContext.HTTPMethod, req.RequestContext.Path)
	fmt.Printf("PathParameters: %+v\n", req.PathParameters)

	// Fix HTTP Method - try multiple sources
	if req.HTTPMethod == "" {
		if req.RequestContext.HTTPMethod != "" {
			req.HTTPMethod = req.RequestContext.HTTPMethod
		} else {
			// Fallback to headers
			if method, exists := req.Headers["X-HTTP-Method-Override"]; exists {
				req.HTTPMethod = method
			} else if method, exists := req.Headers["x-http-method-override"]; exists {
				req.HTTPMethod = method
			} else {
				// Last resort - default to GET for safety
				req.HTTPMethod = "GET"
			}
		}
	}

	// Fix Path
	if req.Path == "" {
		if req.RequestContext.Path != "" {
			req.Path = req.RequestContext.Path
		} else if proxy, exists := req.PathParameters["proxy"]; exists && proxy != "" {
			req.Path = "/" + proxy
		} else {
			req.Path = "/"
		}
	}

	// Remove trailing slash if present (except for root)
	if req.Path != "/" && strings.HasSuffix(req.Path, "/") {
		req.Path = strings.TrimSuffix(req.Path, "/")
	}

	fmt.Printf("Final Request - Method: '%s', Path: '%s'\n", req.HTTPMethod, req.Path)

	return ginLambda.ProxyWithContext(ctx, req)
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
