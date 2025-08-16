package schema

import "time"

type userType string
type userStatus string

type Users struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	Username   string     `json:"userName" gorm:"not null"`
	UserEmail  string     `json:"userEmail" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password   string     `json:"-"` // hide from JSON responses
	UserType   userType   `json:"userType" gorm:"type:varchar(20);check:user_type IN ('staff');not null"`
	UserStatus userStatus `json:"userStatus" gorm:"type:varchar(20);check:user_status IN ('active','inactive');not null;default:'active'"`
}

type Visitor struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
	VisitorName    string     `json:"visitorName" gorm:"not null"`
	VisitorEmail   string     `json:"visitorEmail" gorm:"type:varchar(100);uniqueIndex;"`
	VisitorPhone   string     `json:"visitorPhone" gorm:"type:varchar(15);uniqueIndex;not null"`
	VisitorAddress string     `json:"visitorAddress" gorm:"type:varchar(255)"`
}
