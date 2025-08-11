package schema

import "time"

type userType string
type userStatus string

type Users struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Username   string     `json:"username" gorm:"not null"`
	UserEmail  string     `json:"user_email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password   string     `json:"-"` // hide from JSON responses
	UserType   userType   `json:"user_type" gorm:"type:varchar(20);check:user_type IN ('staff');not null"`
	UserStatus userStatus `json:"user_status" gorm:"type:varchar(20);check:user_status IN ('active','inactive');not null;default:'active'"`
}
