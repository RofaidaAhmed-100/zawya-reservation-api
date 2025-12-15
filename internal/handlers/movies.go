package handlers

import (
	"zawyaReservation/internal/database"
	"zawyaReservation/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateMovieRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes" binding:"required"`
	Genre           string    `json:"genre"`
	Rating          string    `json:"rating"`
	PosterURL       string    `json:"poster_url"`
	ReleaseDate     time.Time `json:"release_date"`
}


func GetMovies(c *gin.Context) {
	var movies []models.Movie
	if err := database.DB.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}


func GetMovie(c *gin.Context) {
	movieID := c.Param("id")

	var movie models.Movie
	if err := database.DB.First(&movie, "id = ?", movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

func CreateMovie(c *gin.Context) {
	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie := models.Movie{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Genre:           req.Genre,
		Rating:          req.Rating,
		PosterURL:       req.PosterURL,
		ReleaseDate:     req.ReleaseDate,
	}

	if err := database.DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Movie created successfully",
		"movie":   movie,
	})
}

func UpdateMovie(c *gin.Context) {
	movieID := c.Param("id")

	var movie models.Movie
	if err := database.DB.First(&movie, "id = ?", movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.DurationMinutes = req.DurationMinutes
	movie.Genre = req.Genre
	movie.Rating = req.Rating
	movie.PosterURL = req.PosterURL
	movie.ReleaseDate = req.ReleaseDate

	if err := database.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie updated successfully",
		"movie":   movie,
	})
}


func DeleteMovie(c *gin.Context) {
	movieID := c.Param("id")

	if err := database.DB.Delete(&models.Movie{}, "id = ?", movieID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
