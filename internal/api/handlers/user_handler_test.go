package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-api-go/internal/models"
	"task-api-go/internal/testutils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	err := testutils.SharedTestDB.Create(&models.User{
		Username:     "John",
		PasswordHash: "Doe",
	}).Error

	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response []models.User
	err = json.Unmarshal(resp.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "John", response[0].Username)
}
