package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Email string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Name      string    `gorm:"not null" json:"name"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}
