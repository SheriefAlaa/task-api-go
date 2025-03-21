package handlers

import (
	"net/http"
	"strconv"

	"task-api-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		err := c.ShouldBindBodyWithJSON(&task)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		db.Create(&task)
		c.JSON(http.StatusCreated, task)
	}
}

func getTaskByID(db *gorm.DB, c *gin.Context) (models.Task, error) {
	id := c.Param("id")
	var task models.Task

	taskID, err := strconv.Atoi(id)
	if err != nil {
		return task, err
	}

	db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).First(&task, taskID)

	if task.ID == 0 {
		return task, err
	}

	return task, nil
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var task models.Task

		taskID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		db.First(&task, taskID)
		if task.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Model(&task).Updates(task)
		c.JSON(http.StatusOK, task)
	}
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		task, err := getTaskByID(db, c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		if task.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		result := db.Delete(&task)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func GetTaskByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		task, err := getTaskByID(db, c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		if task.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

func ListTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		db.Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	}
}
