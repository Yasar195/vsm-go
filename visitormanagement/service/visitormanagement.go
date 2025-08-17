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
	UserId         int64  `json:"userId"`
}

type CreateVisitoryResponse struct {
	Message string `json:"message"`
}

type GetVisitorsResponse struct {
	Visitors []schema.Visitor `json:"visitors"`
	Count    int64            `json:"count"`
}

type GetVisitsResponse struct {
	Visits []schema.Visits `json:"visits"`
	Count  int64           `json:"count"`
}

type GetUserRequest struct {
	PageSize int64
	Page     int64
	Search   string
}

type CreateVisitsInput struct {
	UserId       int64  `json:"userId" validate:"required"`
	VisitorId    int64  `json:"visitorId" validate:"required"`
	VisitPurpose string `json:"visitPurpose"`
}

func CreateVisitor(data CreateVisitorRequest) utility.Response[CreateVisitoryResponse] {
	var visitor = schema.Visitor{
		VisitorName:    data.VisitorName,
		VisitorEmail:   data.VisitorEmail,
		VisitorPhone:   data.VisitorPhone,
		VisitorAddress: data.VisitorAddress,
		CreatedUserID:  uint(data.UserId),
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

func CreateVisits(data CreateVisitsInput) utility.Response[CreateVisitoryResponse] {
	var visit = schema.Visits{
		UserID:        uint(data.UserId),
		VisitorID:     uint(data.VisitorId),
		VisitPurpose:  data.VisitPurpose,
		CreatedUserID: uint(data.UserId),
	}

	err := db.DB.Create(&visit).Error

	if err != nil {
		return utility.Response[CreateVisitoryResponse]{
			Success:    false,
			Message:    "Visit creation failed",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[CreateVisitoryResponse]{
		Success: true,
		Data: &CreateVisitoryResponse{
			Message: "Visit created successfullys",
		},
		Message:    "Visit created successfully",
		StatusCode: http.StatusOK,
	}
}

func GetVisits(data GetUserRequest) utility.Response[GetVisitsResponse] {
	offset := (data.Page - 1) * data.PageSize
	var visits []schema.Visits
	var count int64
	var err, cerr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = db.DB.Model(&schema.Visits{}).Preload("User").Preload("Visitor").Offset(int(offset)).Limit(int(data.PageSize)).Find(&visits).Error
	}()

	go func() {
		defer wg.Done()
		cerr = db.DB.Model(&schema.Visits{}).Count(&count).Error
	}()

	wg.Wait()

	if err != nil {
		return utility.Response[GetVisitsResponse]{
			Success:    false,
			Message:    "Failed to fetch visits",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if cerr != nil {
		return utility.Response[GetVisitsResponse]{
			Success:    false,
			Message:    "Failed to fetch visits",
			Data:       nil,
			Error:      cerr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[GetVisitsResponse]{
		Success: true,
		Data: &GetVisitsResponse{
			Visits: visits,
			Count:  count,
		},
		Message:    "Visits fetched successfully",
		StatusCode: http.StatusOK,
	}
}
