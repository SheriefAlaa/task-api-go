package handlers

import (
	"net/http"
	"task-api-go/internal/api/requests"
	"task-api-go/internal/auth/jwt"
	"task-api-go/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignupUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.SignupInput
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if req.Password != req.PasswordConfirmation {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
			return
		}

		var existingUser models.User
		result := db.Where("username = ?", req.Username).First(&existingUser)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
		}

		result = db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		user.PasswordHash = ""
		c.JSON(http.StatusCreated, user)
	}
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.LoginInput
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user models.User
		result := db.Where("username = ?", req.Username).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := jwt.GenerateUserJwtToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"user_id": user.ID,
		})
	}
}
