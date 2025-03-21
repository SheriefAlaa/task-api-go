package handlers

import (
	"net/http"
	"strconv"
	"task-api-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment creates a new comment for a task
func CreateComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment models.Comment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var task models.Task
		if err := db.First(&task, comment.TaskID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}
		comment.UserID = userID.(uint)

		if err := db.Create(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, comment)
	}
}

// GetCommentsByTaskID returns all comments for a specific task
func GetCommentsByTaskID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		var comments []models.Comment
		if err := db.Where("task_id = ?", taskID).
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, username")
			}).
			Find(&comments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comments)
	}
}

// UpdateComment updates an existing comment
func UpdateComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		var comment models.Comment
		if err := db.First(&comment, commentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}

		// Authorize: only the comment creator can update it.
		//  TODO: use OpenFGA to check if the user is the comment creator
		userID, _ := c.Get("userID")
		if comment.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this comment"})
			return
		}

		// Bind new comment data
		var updateData struct {
			Comment string `json:"comment"`
		}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		comment.Comment = updateData.Comment
		if err := db.Save(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}

// DeleteComment deletes a comment
func DeleteComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		var comment models.Comment
		if err := db.First(&comment, commentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}

		// Authorize: only the comment creator can delete it
		// TODO: use OpenFGA
		userID, _ := c.Get("userID")
		if comment.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this comment"})
			return
		}

		if err := db.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
