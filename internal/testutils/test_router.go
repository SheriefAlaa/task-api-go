package testutils

import (
	"task-api-go/internal/api/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter(db)
}
