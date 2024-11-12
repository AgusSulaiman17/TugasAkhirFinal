package config

import (
    "fmt"
    "log"
    "os"
    "backend/models"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Read database configuration from environment variables
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var dbErr error
    DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if dbErr != nil {
        log.Fatal("Failed to connect to the database:", dbErr)
    }

    // Auto-migrate models to create tables if they don't exist
dbErr = DB.AutoMigrate(&models.Genre{}, &models.Penulis{}, &models.Buku{}, &models.Peminjaman{})
if dbErr != nil {
    log.Fatal("Failed to migrate database:", dbErr)
}


    fmt.Println("Database connected and migrated successfully")
}
