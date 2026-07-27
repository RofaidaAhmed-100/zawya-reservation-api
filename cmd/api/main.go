package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zawyaReservation/internal/database"
	"zawyaReservation/internal/handlers"
	"zawyaReservation/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Println("Warning: JWT_SECRET is not set. Using empty secret is insecure.")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required database environment variables: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME")
	}

	database.Connect()
	database.Migrate()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	authLimiter := middleware.NewRateLimiter(10, time.Minute)

	auth := router.Group("/api/auth")
	auth.Use(authLimiter.Middleware())
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

		api.GET("/movies/:id/showtimes", handlers.GetShowtimesForMovie)

		api.GET("/halls", handlers.GetHalls)
		api.GET("/halls/:id", handlers.GetHall)

		api.GET("/showtimes/:id", handlers.GetShowtime)
		api.GET("/showtimes/:id/seats", handlers.GetAvailableSeats)
	}

	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{

		admin.POST("/movies", handlers.CreateMovie)
		admin.PUT("/movies/:id", handlers.UpdateMovie)
		admin.DELETE("/movies/:id", handlers.DeleteMovie)

		admin.POST("/movies/:id/showtimes", handlers.CreateShowtime)

		admin.POST("/halls", handlers.CreateHall)
		admin.PUT("/halls/:id", handlers.UpdateHall)
		admin.DELETE("/halls/:id", handlers.DeleteHall)
		admin.POST("/halls/:id/seats", handlers.CreateSeatsForHall)

		admin.DELETE("/showtimes/:id", handlers.DeleteShowtime)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	database.Close()

	log.Println("Server exited gracefully")
}
