package schema

import (
	"gorm.io/gorm"
)

type userType string
type userStatus string

type Users struct {
	gorm.Model
	Username   string `gorm:"not null"`
	UserEmail  string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password   string
	UserType   userType   `gorm:"type:varchar(20);check:user_type IN ('staff');not null"`
	UserStatus userStatus `gorm:"type:varchar(20);check:user_status IN ('active','inactive');not null;default:'active'"`
}
