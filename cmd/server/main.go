package main

import (
	"log"
	"net/http"

	"github.com/sheriefalaa/task-api-go/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB connection
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

	// Setup router using our setupRouter function
	router := setupRouter()

	log.Println("Server listening at :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// setupRouter configures the Gin router with API endpoints.
func setupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "API is healthy!"})
		})
	}

	return router
}
