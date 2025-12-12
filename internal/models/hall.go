package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Hall struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	TotalSeats int       `gorm:"not null" json:"total_seats"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (h *Hall) BeforeCreate(tx *gorm.DB) error {
	h.ID = uuid.New().String()
	return nil
}

type Seat struct {
	ID         string `gorm:"type:char(36);primaryKey" json:"id"`
	HallID     string `gorm:"type:char(36);not null" json:"hall_id"`
	RowNumber  int    `gorm:"not null" json:"row_number"`
	SeatNumber int    `gorm:"not null" json:"seat_number"`
	SeatType   string `gorm:"type:varchar(50);default:'regular'" json:"seat_type"`
	
	Hall       Hall   `gorm:"foreignKey:HallID" json:"-"`
}

func (s *Seat) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New().String()
	return nil
}
