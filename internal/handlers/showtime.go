package handlers

import (
	"zawyaReservation/internal/database"
	"zawyaReservation/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateShowtimeRequest struct {
	MovieID   string    `json:"movie_id" binding:"required"`
	HallID    string    `json:"hall_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	BasePrice float64   `json:"base_price" binding:"required"`
}

type SeatAvailability struct {
	Seat      models.Seat `json:"seat"`
	Available bool        `json:"available"`
	Price     float64     `json:"price"`
}


func CreateShowtime(c *gin.Context) {
	var req CreateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	var movie models.Movie
	if err := database.DB.First(&movie, "id = ?", req.MovieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	
	var hall models.Hall
	if err := database.DB.First(&hall, "id = ?", req.HallID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hall not found"})
		return
	}

	endTime := req.StartTime.Add(time.Duration(movie.DurationMinutes) * time.Minute)

	
	var overlapping models.Showtime
	err := database.DB.Where("hall_id = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
		req.HallID, req.StartTime, req.StartTime, endTime, endTime).First(&overlapping).Error
	
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Hall is already booked for this time slot"})
		return
	}

	
	showtime := models.Showtime{
		MovieID:   req.MovieID,
		HallID:    req.HallID,
		StartTime: req.StartTime,
		EndTime:   endTime,
		BasePrice: req.BasePrice,
	}

	if err := database.DB.Create(&showtime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create showtime"})
		return
	}

	database.DB.Preload("Movie").Preload("Hall").First(&showtime, "id = ?", showtime.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Showtime created successfully",
		"showtime": showtime,
	})
}

func GetShowtimesForMovie(c *gin.Context) {
	movieID := c.Param("id")

	
	var movie models.Movie
	if err := database.DB.First(&movie, "id = ?", movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	
	var showtimes []models.Showtime
	database.DB.Where("movie_id = ? AND start_time > ?", movieID, time.Now()).
		Preload("Hall").
		Order("start_time").
		Find(&showtimes)

	c.JSON(http.StatusOK, gin.H{
		"movie":     movie,
		"showtimes": showtimes,
	})
}


func GetShowtime(c *gin.Context) {
	showtimeID := c.Param("id")

	var showtime models.Showtime
	if err := database.DB.Preload("Movie").Preload("Hall").First(&showtime, "id = ?", showtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtime": showtime})
}

func GetAvailableSeats(c *gin.Context) {
	showtimeID := c.Param("id")

	
	var showtime models.Showtime
	if err := database.DB.First(&showtime, "id = ?", showtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	
	var seats []models.Seat
	database.DB.Where("hall_id = ?", showtime.HallID).Order("row_number, seat_number").Find(&seats)

	
	var availability []SeatAvailability
	for _, seat := range seats {
		price := calculateSeatPrice(showtime.BasePrice, seat.SeatType)
		availability = append(availability, SeatAvailability{
			Seat:      seat,
			Available: true,
			Price:     price,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"showtime":     showtime,
		"availability": availability,
	})
}


func calculateSeatPrice(basePrice float64, seatType string) float64 {
	switch seatType {
	case "premium":
		return basePrice * 1.5 // 50% more
	case "vip":
		return basePrice * 2.0 // 100% more (double)
	default:
		return basePrice
	}
}


func DeleteShowtime(c *gin.Context) {
	showtimeID := c.Param("id")

	if err := database.DB.Delete(&models.Showtime{}, "id = ?", showtimeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete showtime"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Showtime deleted successfully"})
}