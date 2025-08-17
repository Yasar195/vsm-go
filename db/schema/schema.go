package schema

import "time"

type UserType string
type UserStatus string
type VisitorStatus string

type Users struct {
	ID         uint       `json:"id" gorm:"column:id"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at"`
	Username   string     `json:"userName" gorm:"column:user_name;not null"`
	UserEmail  string     `json:"userEmail" gorm:"column:user_email;type:varchar(100);uniqueIndex;not null"`
	Password   string     `json:"-" gorm:"column:password"` // hide from JSON responses
	UserType   UserType   `json:"userType" gorm:"column:user_type;type:varchar(20);check:user_type IN ('staff', 'host');not null"`
	UserStatus UserStatus `json:"userStatus" gorm:"column:user_status;type:varchar(20);check:user_status IN ('active','inactive');not null;default:'active'"`
}

type Visitor struct {
	ID             uint       `json:"id" gorm:"column:id"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at"`
	VisitorName    string     `json:"visitorName" gorm:"column:visitor_name;not null"`
	VisitorEmail   string     `json:"visitorEmail" gorm:"column:visitor_email;type:varchar(100);uniqueIndex;"`
	VisitorPhone   string     `json:"visitorPhone" gorm:"column:visitor_phone;type:varchar(15);uniqueIndex;not null"`
	VisitorAddress string     `json:"visitorAddress" gorm:"column:visitor_address;type:varchar(255)"`
	IsVerified     bool       `json:"isVerified" gorm:"column:is_verified;default:false"`
	CreatedUserID  uint       `json:"createdUserId" gorm:"column:created_user_id;not null"`
	CreatedUser    Users      `json:"createdUser" gorm:"foreignKey:CreatedUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Visits struct {
	ID            uint          `json:"id" gorm:"column:id"`
	CreatedAt     time.Time     `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time     `json:"updatedAt" gorm:"column:updated_at"`
	UserID        uint          `json:"userId" gorm:"column:user_id;not null"`
	VisitorID     uint          `json:"visitorId" gorm:"column:visitor_id;not null"`
	User          Users         `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Visitor       Visitor       `json:"visitor" gorm:"foreignKey:VisitorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	VisitStatus   VisitorStatus `json:"visitStatus" gorm:"column:visit_status;type:varchar(20);check:visit_status IN ('waiting', 'meeting', 'canceled', 'completed');not null;default:'waiting'"`
	VisitPurpose  string        `json:"visitPurpose" gorm:"column:visit_purpose;type:varchar(255)"`
	CreatedUserID uint          `json:"createdUserId" gorm:"column:created_user_id;not null"`
	CreatedUser   Users         `json:"createdUser" gorm:"foreignKey:CreatedUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Notifications struct {
	ID        uint      `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UserID    uint      `json:"userId" gorm:"column:user_id;not null"`
	Title     string    `json:"title" gorm:"column:title;type:varchar(255);not null"`
	Message   string    `json:"message" gorm:"column:message;type:varchar(255);not null"`
	User      Users     `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsRead    bool      `json:"isRead" gorm:"column:is_read;default:false"`
}

type Logs struct {
	ID        uint      `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UserID    uint      `json:"userId" gorm:"column:user_id;not null"`
	Action    string    `json:"action" gorm:"column:action;type:varchar(255);not null"`
}
