package main

import (
	"log"
	"task-api-go/internal/api/routes"
	"task-api-go/internal/utils"
)

func main() {
	db, err := utils.InitDB()
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	// Retrieve the underlying *sql.DB and defer its close
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB from GORM:", err)
	}
	defer sqlDB.Close()

	router := routes.SetupRouter(db)

	log.Println("Server listening at :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
