package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Movie struct {
	ID              string    `gorm:"type:char(36);primaryKey" json:"id"`
	Title           string    `gorm:"not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	DurationMinutes int       `gorm:"not null" json:"duration_minutes"`
	Genre           string    `json:"genre"`
	Rating          string    `json:"rating"`
	PosterURL       string    `json:"poster_url"`
	ReleaseDate     time.Time `json:"release_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New().String()
	return nil
}
