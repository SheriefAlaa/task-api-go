// internal/api/routes/setup.go
package routes

import (
	"net/http"
	"task-api-go/internal/api/handlers"
	"task-api-go/internal/api/middleware"

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

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handlers.SignupUser(db))
			auth.POST("/login", handlers.LoginUser(db))
		}

		// Protected routes - require authentication
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// Users
			protected.GET("/users", handlers.GetAllUsers(db))
			protected.GET("/users/:id/tasks", handlers.GetUserTasks(db))

			// Tasks
			protected.POST("/tasks", handlers.CreateTask(db))
			protected.GET("/tasks", handlers.ListTasks(db))
			protected.GET("/tasks/:id", handlers.GetTaskByID(db))
			protected.PUT("/tasks/:id", handlers.UpdateTask(db))
			protected.DELETE("/tasks/:id", handlers.DeleteTask(db))
		}
	}

	return router
}
