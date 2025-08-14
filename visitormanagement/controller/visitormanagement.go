package visitormanagementcontrollers

import (
	visitormanagementservice "visitor-management-system/visitormanagement/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func CreateVisitor(c *gin.Context) {
	var body visitormanagementservice.CreateVisitorRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"success":    false,
			"message":    "visitor creation failed",
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

	resp := visitormanagementservice.CreateVisitor(body)
	c.JSON(resp.StatusCode, resp)
}
