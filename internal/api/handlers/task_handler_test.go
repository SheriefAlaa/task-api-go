package handlers_test

import (
	"testing"

	"task-api-go/internal/api/handlers"
	"task-api-go/internal/testutils"

	"github.com/gin-gonic/gin"
)

func TestCreateTask(t *testing.T) {
	r := gin.Default()
	t.Parallel()
	testutils.CascadeDB(t)

	// r := testutils.SetupTestRouter(testutils.SharedTestDB)

	r.POST("/api/v1/tasks", handlers.CreateTask(testutils.SharedTestDB))

}

// TODO: Add more tests
