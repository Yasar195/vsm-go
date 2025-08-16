package usermanagementcontroller

import (
	"strconv"
	usermanagementservice "visitor-management-system/iam/usermanagement/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

func GetUsers(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	search := c.DefaultQuery("search", "")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 10
	}

	resp := usermanagementservice.GetUsers(usermanagementservice.GetUserRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	c.JSON(resp.StatusCode, resp)
}

var validate = validator.New()

func CreateUser(c *gin.Context) {
	var body usermanagementservice.CreateUserInput
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"success":    false,
			"message":    "User creation failed",
			"error":      err.Error(),
			"data":       nil,
			"statusCode": 400,
		})
		return
	}

	if err := validate.Struct(&body); err != nil {
		c.JSON(400, gin.H{
			"success":    false,
			"message":    "Validation failed",
			"error":      err.Error(),
			"data":       nil,
			"statusCode": 400,
		})
		return
	}

	resp := usermanagementservice.CreateUser(body)
	c.JSON(resp.StatusCode, resp)
}
