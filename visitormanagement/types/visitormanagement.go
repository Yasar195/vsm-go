package visitormanagementtypes

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
	Visitors []VisitorOriginalResponse `json:"visitors"`
	Count    int64                     `json:"count"`
}

type VisitResponse struct {
	ID           uint   `json:"id"`
	VisitStatus  string `json:"visitStatus"`
	VisitPurpose string `json:"visitPurpose"`

	UserID    uint `json:"-" gorm:"column:user_id"`
	VisitorID uint `json:"-" gorm:"column:visitor_id"`

	User    UserResponse    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Visitor VisitorResponse `json:"visitor" gorm:"foreignKey:VisitorID;references:ID"`
}

type VisitorOriginalResponse struct {
	ID             uint   `json:"id"`
	VisitorName    string `json:"visitorName"`
	VisitorEmail   string `json:"visitorEmail"`
	VisitorPhone   string `json:"visitorPhone"`
	VisitorAddress string `json:"visitorAddress"`
	IsVerified     bool   `json:"isVerified"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

type VisitorResponse struct {
	ID           uint   `json:"id"`
	VisitorName  string `json:"visitorName"`
	VisitorEmail string `json:"visitorEmail"`
	VisitorPhone string `json:"visitorPhone"`
}

func (UserResponse) TableName() string {
	return "users"
}

func (VisitorResponse) TableName() string {
	return "visitors"
}

type GetVisitsResponse struct {
	Visits []VisitResponse `json:"visits"`
	Count  int64           `json:"count"`
}

type GetUserRequest struct {
	PageSize      int64
	Page          int64
	Search        string
	VisitorStatus *string
	VisitorId     *int64
}

type CreateVisitsInput struct {
	UserId       int64  `json:"userId" validate:"required"`
	VisitorId    int64  `json:"visitorId" validate:"required"`
	VisitPurpose string `json:"visitPurpose"`
}
