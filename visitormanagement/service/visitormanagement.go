package visitormanagementservice

import (
	"net/http"
	"sync"
	"visitor-management-system/db"
	"visitor-management-system/db/schema"
	"visitor-management-system/utility"
	visitormanagementtypes "visitor-management-system/visitormanagement/types"

	"gorm.io/gorm"
)

func CreateVisitor(data visitormanagementtypes.CreateVisitorRequest) utility.Response[visitormanagementtypes.CreateVisitoryResponse] {
	var visitor = schema.Visitor{
		VisitorName:    data.VisitorName,
		VisitorEmail:   data.VisitorEmail,
		VisitorPhone:   data.VisitorPhone,
		VisitorAddress: data.VisitorAddress,
		CreatedUserID:  uint(data.UserId),
	}

	err := db.DB.Create(&visitor).Error

	if err != nil {
		return utility.Response[visitormanagementtypes.CreateVisitoryResponse]{
			Success:    false,
			Message:    "User creation failed",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[visitormanagementtypes.CreateVisitoryResponse]{
		Success: true,
		Data: &visitormanagementtypes.CreateVisitoryResponse{
			Message: "Visitor created successfully",
		},
		Message:    "Visitor creation success",
		StatusCode: http.StatusOK,
	}
}

func GetVisitors(data visitormanagementtypes.GetUserRequest) utility.Response[visitormanagementtypes.GetVisitorsResponse] {

	offset := (data.Page - 1) * data.PageSize
	var visitors []visitormanagementtypes.VisitorOriginalResponse
	var count int64
	var err, cerr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = db.DB.Model(&schema.Visitor{}).Select("id, visitor_name, visitor_email, visitor_phone, visitor_address, is_verified").Where("visitor_name ILIKE ?", "%"+data.Search+"%").Offset(int(offset)).Limit(int(data.PageSize)).Find(&visitors).Error
	}()

	go func() {
		defer wg.Done()
		cerr = db.DB.Model(&schema.Visitor{}).Where("visitor_name ILIKE ?", "%"+data.Search+"%").Count(&count).Error
	}()

	wg.Wait()

	if err != nil {
		return utility.Response[visitormanagementtypes.GetVisitorsResponse]{
			Success:    false,
			Message:    "Failed to fetch visitors",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if cerr != nil {
		return utility.Response[visitormanagementtypes.GetVisitorsResponse]{
			Success:    false,
			Message:    "Failed to fetch visitors",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[visitormanagementtypes.GetVisitorsResponse]{
		Success: true,
		Data: &visitormanagementtypes.GetVisitorsResponse{
			Visitors: visitors,
			Count:    count,
		},
		Message:    "Visitors fetched successfully",
		StatusCode: http.StatusOK,
	}
}

func CreateVisits(data visitormanagementtypes.CreateVisitsInput) utility.Response[visitormanagementtypes.CreateVisitoryResponse] {
	var visit = schema.Visits{
		UserID:        uint(data.UserId),
		VisitorID:     uint(data.VisitorId),
		VisitPurpose:  data.VisitPurpose,
		CreatedUserID: uint(data.UserId),
	}

	err := db.DB.Create(&visit).Error

	if err != nil {
		return utility.Response[visitormanagementtypes.CreateVisitoryResponse]{
			Success:    false,
			Message:    "Visit creation failed",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[visitormanagementtypes.CreateVisitoryResponse]{
		Success: true,
		Data: &visitormanagementtypes.CreateVisitoryResponse{
			Message: "Visit created successfullys",
		},
		Message:    "Visit created successfully",
		StatusCode: http.StatusOK,
	}
}

func GetVisits(data visitormanagementtypes.GetUserRequest) utility.Response[visitormanagementtypes.GetVisitsResponse] {
	offset := (data.Page - 1) * data.PageSize
	var visits []visitormanagementtypes.VisitResponse
	var count int64
	var err, cerr error
	var wg sync.WaitGroup

	query := db.DB.Model(&schema.Visits{}).Select("id, visit_status, visit_purpose, user_id, visitor_id")

	if data.VisitorStatus != nil && *data.VisitorStatus != "" {
		query = query.Where("visit_status = ?", data.VisitorStatus)
	}
	if data.VisitorId != nil {
		query = query.Where("visitor_id = ?", *data.VisitorId)
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		err = query.Session(&gorm.Session{}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, user_name, user_email")
		}).Preload("Visitor", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, visitor_name, visitor_email, visitor_phone")
		}).Offset(int(offset)).Limit(int(data.PageSize)).Find(&visits).Error
	}()

	go func() {
		defer wg.Done()
		cerr = query.Session(&gorm.Session{}).Count(&count).Error
	}()

	wg.Wait()

	if err != nil {
		return utility.Response[visitormanagementtypes.GetVisitsResponse]{
			Success:    false,
			Message:    "Failed to fetch visits",
			Data:       nil,
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if cerr != nil {
		return utility.Response[visitormanagementtypes.GetVisitsResponse]{
			Success:    false,
			Message:    "Failed to fetch visits",
			Data:       nil,
			Error:      cerr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return utility.Response[visitormanagementtypes.GetVisitsResponse]{
		Success: true,
		Data: &visitormanagementtypes.GetVisitsResponse{
			Visits: visits,
			Count:  count,
		},
		Message:    "Visits fetched successfully",
		StatusCode: http.StatusOK,
	}
}
