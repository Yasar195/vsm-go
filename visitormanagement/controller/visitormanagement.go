package visitormanagementcontrollers

import (
	"net/http"
	"strconv"
	visitormanagementservice "visitor-management-system/visitormanagement/service"
	visitormanagementtypes "visitor-management-system/visitormanagement/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var validate = validator.New()

func CreateVisitor(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	var body visitormanagementtypes.CreateVisitorRequest
	body.UserId = userID
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

func GetVisitors(c *gin.Context) {
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

	resp := visitormanagementservice.GetVisitors(visitormanagementtypes.GetUserRequest{
		PageSize: pageSize,
		Page:     page,
		Search:   search,
	})
	c.JSON(resp.StatusCode, resp)
}

func CreateVisits(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	var body visitormanagementtypes.CreateVisitsInput
	body.UserId = userID
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"success":    false,
			"message":    "visit creation failed",
			"error":      err.Error(),
			"data":       nil,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	if err := validate.Struct(&body); err != nil {
		c.JSON(400, gin.H{
			"success":    false,
			"message":    "Validation failed",
			"error":      err.Error(),
			"data":       nil,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	resp := visitormanagementservice.CreateVisits(body)
	c.JSON(resp.StatusCode, resp)
}

func GetVisits(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	search := c.DefaultQuery("search", "")
	visitorStatusStr := c.Query("visitorStatus")
	visitorIdStr := c.Query("visitorId")

	var visitorId *int64
	if visitorIdStr != "" {
		if parsedId, err := strconv.ParseInt(visitorIdStr, 10, 64); err == nil {
			visitorId = &parsedId
		}
	}

	var visitorStatus *string
	if visitorStatusStr != "" {
		visitorStatus = &visitorStatusStr
	}

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 10
	}

	resp := visitormanagementservice.GetVisits(visitormanagementtypes.GetUserRequest{
		PageSize:      pageSize,
		Page:          page,
		Search:        search,
		VisitorStatus: visitorStatus,
		VisitorId:     visitorId,
	})
	c.JSON(resp.StatusCode, resp)
}
