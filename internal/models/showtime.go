package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Showtime struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	MovieID   string    `gorm:"type:char(36);not null" json:"movie_id"`
	HallID    string    `gorm:"type:char(36);not null" json:"hall_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	BasePrice float64   `gorm:"type:decimal(10,2);not null" json:"base_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Movie Movie `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Hall  Hall  `gorm:"foreignKey:HallID" json:"hall,omitempty"`
}

func (s *Showtime) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New().String()
	return nil
}
