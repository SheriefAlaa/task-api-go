package handlers

import (
	"net/http"

	"task-api-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		db.Preload("Tasks").Find(&users)

		c.JSON(http.StatusOK, users)
	}
}

func GetUserTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		id := c.Param("id")

		db.Where("assignee_id = ?", id).Find(&tasks)

		c.JSON(http.StatusOK, tasks)
	}
}
