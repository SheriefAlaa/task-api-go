package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"task-api-go/internal/auth/jwt"
	"task-api-go/internal/models"
	"task-api-go/internal/testutils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	user := models.User{
		Username:     "John",
		PasswordHash: "Doe",
	}

	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	token, err := jwt.GenerateUserJwtToken(user.ID)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response []models.User
	err = json.Unmarshal(resp.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "John", response[0].Username)
}

func TestGetUserTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	user1 := models.User{
		Username:     "John",
		PasswordHash: "Doe",
	}
	user2 := models.User{
		Username:     "Jane",
		PasswordHash: "Doe",
	}

	err := testutils.SharedTestDB.Create(&user1).Error
	assert.NoError(t, err)
	err = testutils.SharedTestDB.Create(&user2).Error
	assert.NoError(t, err)

	// Now create tasks with valid user IDs
	task1 := models.Task{
		Title:       "Task 1",
		Description: "Task 1 content",
		AssigneeID:  user1.ID,
	}
	task2 := models.Task{
		Title:       "Task 2",
		Description: "Task 2 content",
		AssigneeID:  user2.ID,
		Status:      models.TaskStatus(models.StatusInProgress),
	}

	err = testutils.SharedTestDB.Create(&task1).Error
	assert.NoError(t, err)
	err = testutils.SharedTestDB.Create(&task2).Error
	assert.NoError(t, err)
	token, err := jwt.GenerateUserJwtToken(user1.ID)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/users/"+strconv.Itoa(int(user2.ID))+"/tasks", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response []models.Task
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Task 2", response[0].Title)
	assert.Equal(t, models.TaskStatus("in_progress"), response[0].Status)
}
