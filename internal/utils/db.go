package utils

import (
	"fmt"
	"log"
	"os"
	"sync"

	"task-api-go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"))

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to the database:", err)
		}

		err = db.AutoMigrate(&models.User{}, &models.Task{}, &models.Comment{})
		if err != nil {
			log.Fatal("Failed to run migrations:", err)
		}

		log.Println("Successfully connected to DB and ran migrations!")
	})

	return db, err
}
