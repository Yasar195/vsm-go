package logmanagementcontroller

import (
	"strconv"
	logmanagementservice "visitor-management-system/logmanagement/service"

	"github.com/gin-gonic/gin"
)

func GetLogs(c *gin.Context) {
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

	resp := logmanagementservice.GetApplicationLogs(logmanagementservice.GetApplicationLogsRequest{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	})
	c.JSON(resp.StatusCode, resp)
}
