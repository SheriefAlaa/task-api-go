package handlers

import (
	"net/http"

	"task-api-go/internal/api/requests"
	"task-api-go/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignupUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input requests.SignupInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		if input.Password != input.PasswordConfirmation {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password and password confirmation do not match"})
			return
		}

		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			Username:     input.Username,
			PasswordHash: string(hashedPwd),
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":       user.ID,
			"username": user.Username,
		})
	}
}
