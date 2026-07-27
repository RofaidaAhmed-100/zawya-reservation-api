package database

import (
	"fmt"
	"log"
	"os"

	"zawyaReservation/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: check your DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, and DB_NAME")
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Hall{},
		&models.Seat{},
		&models.Showtime{},
		&models.Reservation{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed")
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Failed to get underlying DB:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
		return
	}
	log.Println("Database connection closed")
}
