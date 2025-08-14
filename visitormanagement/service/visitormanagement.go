package visitormanagementservice

import (
	"net/http"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"
)

type CreateVisitorRequest struct {
	VisitorName    string `json:"visitorName" validate:"required"`
	VisitorEmail   string `json:"visitorEmail" validate:"required,email"`
	VisitorPhone   string `json:"visitorPhone" validate:"required"`
	VisitorAddress string `json:"visitorAddress" validate:"required"`
}

type CreateVisitoryResponse struct {
	Message string `json:"message"`
}

func CreateVisitor(data CreateVisitorRequest) utility.Response[CreateVisitoryResponse] {
	var visitor = schema.Visitor{
		VisitorName:    data.VisitorName,
		VisitorEmail:   data.VisitorEmail,
		VisitorPhone:   data.VisitorPhone,
		VisitorAddress: data.VisitorAddress,
	}

	err := db.DB.Create(&visitor).Error

	if err != nil {
		return utility.Response[CreateVisitoryResponse]{
			Success:    false,
			Message:    "Failed to create visitor",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[CreateVisitoryResponse]{
		Success: true,
		Data: &CreateVisitoryResponse{
			Message: "Visitor created successfully",
		},
		Message:    "Visitor creation success",
		StatusCode: http.StatusOK,
	}
}
