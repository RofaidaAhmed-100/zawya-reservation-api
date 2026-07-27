package handlers

import (
	"net/http"
	"strings"
	"time"
	"zawyaReservation/internal/database"
	"zawyaReservation/internal/models"
	"zawyaReservation/internal/money"

	"github.com/gin-gonic/gin"
)

type CreateShowtimeRequest struct {
	MovieID   string    `json:"movie_id" binding:"required"`
	HallID    string    `json:"hall_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	BasePrice float64   `json:"base_price" binding:"required"`
	Currency  string    `json:"currency" binding:"required"`
}

type SeatAvailability struct {
	Seat       models.Seat `json:"seat"`
	Available  bool        `json:"available"`
	Price      uint64      `json:"price"`
	PriceFloat float64     `json:"price_float"`
	Currency   string      `json:"currency"`
}

func CreateShowtime(c *gin.Context) {
	var req CreateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.MovieID = strings.TrimSpace(req.MovieID)
	req.HallID = strings.TrimSpace(req.HallID)
	req.Currency = strings.TrimSpace(req.Currency)

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

	currency := req.Currency
	if currency == "" {
		currency = "EGP"
	}

	moneyValue, err := money.NewFromFloat(req.BasePrice, money.Currency(currency))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price or currency"})
		return
	}

	showtime := models.Showtime{
		MovieID:   req.MovieID,
		HallID:    req.HallID,
		StartTime: req.StartTime,
		EndTime:   endTime,
		BasePrice: moneyValue.Amount,
		Currency:  string(moneyValue.Currency),
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

	var reservedSeatIDs []string
	database.DB.Table("reservation_seats").
		Joins("JOIN reservations ON reservations.id = reservation_seats.reservation_id").
		Where("reservations.showtime_id = ? AND reservations.status = ?", showtimeID, "confirmed").
		Pluck("reservation_seats.seat_id", &reservedSeatIDs)

	reservedSet := make(map[string]bool, len(reservedSeatIDs))
	for _, id := range reservedSeatIDs {
		reservedSet[id] = true
	}

	for _, seat := range seats {
		priceMoney := calculateSeatPrice(showtime.BasePrice, showtime.Currency, seat.SeatType)
		available := !reservedSet[seat.ID]
		availability = append(availability, SeatAvailability{
			Seat:       seat,
			Available:  available,
			Price:      priceMoney.Amount,
			PriceFloat: priceMoney.ToFloat(),
			Currency:   string(priceMoney.Currency),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"showtime":     showtime,
		"availability": availability,
	})
}

func calculateSeatPrice(baseAmount uint64, currency string, seatType string) *money.Money {
	baseMoney := money.New(baseAmount, money.Currency(currency))

	var result *money.Money
	var err error

	switch seatType {
	case "premium":
		result, err = baseMoney.MultiplyFloat(1.5)
	case "vip":
		result, err = baseMoney.MultiplyFloat(2.0)
	default:
		result = baseMoney
	}

	if err != nil {
		return baseMoney
	}

	return result
}

func DeleteShowtime(c *gin.Context) {
	showtimeID := c.Param("id")

	if err := database.DB.Delete(&models.Showtime{}, "id = ?", showtimeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete showtime"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Showtime deleted successfully"})
}
