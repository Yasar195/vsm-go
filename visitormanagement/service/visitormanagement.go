package visitormanagementservice

import (
	"net/http"
	"sync"
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

type GetVisitorsResponse struct {
	Visitors []schema.Visitor `json:"visitors"`
	Count    int64            `json:"count"`
}

type GetUserRequest struct {
	PageSize int64
	Page     int64
	Search   string
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
			Message:    "User creation failed",
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

func GetVisitors(data GetUserRequest) utility.Response[GetVisitorsResponse] {

	offset := (data.Page - 1) * data.PageSize
	var visitors []schema.Visitor
	var count int64
	var err, cerr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = db.DB.Model(&schema.Visitor{}).Where("visitor_name ILIKE ?", "%"+data.Search+"%").Offset(int(offset)).Limit(int(data.PageSize)).Find(&visitors).Error
	}()

	go func() {
		defer wg.Done()
		cerr = db.DB.Model(&schema.Visitor{}).Where("visitor_name ILIKE ?", "%"+data.Search+"%").Count(&count).Error
	}()

	wg.Wait()

	if err != nil {
		return utility.Response[GetVisitorsResponse]{
			Success:    false,
			Message:    "Failed to fetch visitors",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if cerr != nil {
		return utility.Response[GetVisitorsResponse]{
			Success:    false,
			Message:    "Failed to fetch visitors",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[GetVisitorsResponse]{
		Success: true,
		Data: &GetVisitorsResponse{
			Visitors: visitors,
			Count:    count,
		},
		Message:    "Visitors fetched successfully",
		StatusCode: http.StatusOK,
	}
}
