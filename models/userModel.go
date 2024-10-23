package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName string    `gorm:"type:varchar(255);not null" json:"firstName" validate:"required,min=2,max=100"`
	LastName  string    `gorm:"type:varchar(255);not null" json:"lastName" validate:"required,min=2,max=100"`
	Email     string    `gorm:"type:varchar(255);not null;unique" json:"email" validate:"required,email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=8"`
	Phone     string    `gorm:"type:varchar(255);not null;unique" json:"phone"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Balance   int64     `gorm:"type:INT;default:0;not null" json:"balance"`
	UserType  string    `gorm:"type:varchar(255);not null;enum('ADMIN','CLIENT');default:'CLIENT'" json:"userType" validate:"required,eq=ADMIN|eq=CLIENT"`
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
