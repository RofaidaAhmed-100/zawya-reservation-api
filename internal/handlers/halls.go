package handlers

import (
	"zawyaReservation/internal/database"
	"zawyaReservation/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateHallRequest struct {
	Name       string `json:"name" binding:"required"`
	TotalSeats int    `json:"total_seats" binding:"required"`
}

type CreateSeatsRequest struct {
	Rows          int   `json:"rows" binding:"required"`
	SeatsPerRow   int   `json:"seats_per_row" binding:"required"`
	PremiumRows   []int `json:"premium_rows"`
	VIPRows       []int `json:"vip_rows"`
}


func GetHalls(c *gin.Context) {
	var halls []models.Hall
	if err := database.DB.Find(&halls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch halls"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"halls": halls})
}

func GetHall(c *gin.Context) {
	hallID := c.Param("id")

	var hall models.Hall
	if err := database.DB.First(&hall, "id = ?", hallID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hall not found"})
		return
	}

	var seats []models.Seat
	database.DB.Where("hall_id = ?", hallID).Order("row_number, seat_number").Find(&seats)

	c.JSON(http.StatusOK, gin.H{
		"hall":  hall,
		"seats": seats,
	})
}

func CreateHall(c *gin.Context) {
	var req CreateHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hall := models.Hall{
		Name:       req.Name,
		TotalSeats: req.TotalSeats,
	}

	if err := database.DB.Create(&hall).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hall"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hall created successfully",
		"hall":    hall,
	})
}


func CreateSeatsForHall(c *gin.Context) {
	hallID := c.Param("id")

	var hall models.Hall
	if err := database.DB.First(&hall, "id = ?", hallID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hall not found"})
		return
	}

	var req CreateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPremium := func(row int) bool {
		for _, r := range req.PremiumRows {
			if r == row {
				return true
			}
		}
		return false
	}

	isVIP := func(row int) bool {
		for _, r := range req.VIPRows {
			if r == row {
				return true
			}
		}
		return false
	}

	
	var seats []models.Seat
	for row := 1; row <= req.Rows; row++ {
		for seat := 1; seat <= req.SeatsPerRow; seat++ {
			seatType := "regular"
			if isVIP(row) {
				seatType = "vip"
			} else if isPremium(row) {
				seatType = "premium"
			}

			seats = append(seats, models.Seat{
				HallID:     hallID,
				RowNumber:  row,
				SeatNumber: seat,
				SeatType:   seatType,
			})
		}
	}

	
	if err := database.DB.Create(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seats"})
		return
	}

	
	hall.TotalSeats = len(seats)
	database.DB.Save(&hall)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seats created successfully",
		"count":   len(seats),
	})
}
