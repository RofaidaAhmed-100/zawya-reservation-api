package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     string    `gorm:"type:char(36);not null" json:"user_id"`
	ShowtimeID string    `gorm:"type:char(36);not null" json:"showtime_id"`
	TotalPrice uint64    `gorm:"not null" json:"total_price"`
	Currency   string    `gorm:"type:char(3);not null" json:"currency"`
	Status     string    `gorm:"type:varchar(20);default:'confirmed'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Showtime Showtime `gorm:"foreignKey:ShowtimeID" json:"showtime,omitempty"`
	Seats    []Seat   `gorm:"many2many:reservation_seats" json:"seats,omitempty"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New().String()
	return nil
}
