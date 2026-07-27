package handlers

import (
	"net/http"
	"strings"

	"zawyaReservation/internal/database"
	"zawyaReservation/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateReservationRequest struct {
	ShowtimeID string   `json:"showtime_id" binding:"required"`
	SeatIDs    []string `json:"seat_ids" binding:"required,min=1"`
}

func CreateReservation(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ShowtimeID = strings.TrimSpace(req.ShowtimeID)
	for i := range req.SeatIDs {
		req.SeatIDs[i] = strings.TrimSpace(req.SeatIDs[i])
	}

	var showtime models.Showtime
	if err := database.DB.First(&showtime, "id = ?", req.ShowtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	var seats []models.Seat
	if err := database.DB.Where("id IN ? AND hall_id = ?", req.SeatIDs, showtime.HallID).Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}
	if len(seats) != len(req.SeatIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more seats not found or not in this hall"})
		return
	}

	tx := database.DB.Begin()

	var reservedCount int64
	if err := tx.Table("reservation_seats").
		Joins("JOIN reservations ON reservations.id = reservation_seats.reservation_id").
		Where("reservations.showtime_id = ? AND reservations.status = ? AND reservation_seats.seat_id IN ?",
			req.ShowtimeID, "confirmed", req.SeatIDs).
		Count(&reservedCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check seat availability"})
		return
	}
	if reservedCount > 0 {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "One or more seats are already reserved"})
		return
	}

	var totalAmount uint64
	for _, seat := range seats {
		priceMoney := calculateSeatPrice(showtime.BasePrice, showtime.Currency, seat.SeatType)
		totalAmount += priceMoney.Amount
	}

	reservation := models.Reservation{
		UserID:     userID,
		ShowtimeID: req.ShowtimeID,
		TotalPrice: totalAmount,
		Currency:   showtime.Currency,
		Status:     "confirmed",
	}

	if err := tx.Create(&reservation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	if err := tx.Model(&reservation).Association("Seats").Append(seats); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate seats"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete reservation"})
		return
	}

	database.DB.Preload("Seats").Preload("Showtime").Preload("Showtime.Movie").Preload("Showtime.Hall").First(&reservation, "id = ?", reservation.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Reservation created successfully",
		"reservation": reservation,
	})
}

func GetUserReservations(c *gin.Context) {
	userID := c.GetString("user_id")

	var reservations []models.Reservation
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Seats").
		Preload("Showtime").
		Preload("Showtime.Movie").
		Preload("Showtime.Hall").
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func GetReservation(c *gin.Context) {
	reservationID := c.Param("id")
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	query := database.DB.Where("id = ?", reservationID)
	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	var reservation models.Reservation
	if err := query.
		Preload("Seats").
		Preload("Showtime").
		Preload("Showtime.Movie").
		Preload("Showtime.Hall").
		First(&reservation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservation": reservation})
}

func CancelReservation(c *gin.Context) {
	reservationID := c.Param("id")
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	query := database.DB.Where("id = ? AND status = ?", reservationID, "confirmed")
	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	var reservation models.Reservation
	if err := query.First(&reservation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found or already cancelled"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservation"})
		return
	}

	reservation.Status = "cancelled"
	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled successfully"})
}
