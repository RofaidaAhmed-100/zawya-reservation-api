package main

import (
	"log"
	"zawyaReservation/internal/database"
	"zawyaReservation/internal/handlers"
	"zawyaReservation/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	
	database.Connect()
	database.Migrate()
	router := gin.Default()

	
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	
	api := router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		
		api.GET("/profile", handlers.GetProfile)

		api.GET("/movies", handlers.GetMovies)
		api.GET("/movies/:id", handlers.GetMovie)

		api.GET("/halls", handlers.GetHalls)
		api.GET("/halls/:id", handlers.GetHall)
	}

	
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		
		admin.POST("/movies", handlers.CreateMovie)
		admin.PUT("/movies/:id", handlers.UpdateMovie)
		admin.DELETE("/movies/:id", handlers.DeleteMovie)

	
		admin.POST("/halls", handlers.CreateHall)
		admin.POST("/halls/:id/seats", handlers.CreateSeatsForHall)
	}

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
