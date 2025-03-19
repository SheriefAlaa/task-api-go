// internal/api/routes/setup.go
package routes

import (
	"net/http"
	"task-api-go/internal/api/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "API is healthy!"})
		})

		// Auth
		v1.POST("/auth/signup", handlers.SignupUser(db))

		// Users
		v1.GET("/users", handlers.GetAllUsers(db))

		// Tasks
		v1.POST("/tasks", handlers.CreateTask(db))
		v1.GET("/tasks", handlers.ListTasks(db))
		v1.GET("/tasks/:id", handlers.GetTaskByID(db))
		v1.PUT("/tasks/:id", handlers.UpdateTask(db))
		v1.DELETE("/tasks/:id", handlers.DeleteTask(db))
	}

	return router
}
