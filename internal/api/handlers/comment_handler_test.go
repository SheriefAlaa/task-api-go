package handlers_test

import (
	"bytes"
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

func TestCreateComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	user := models.User{
		Username:     "commentuser",
		PasswordHash: "password",
	}
	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	task := models.Task{
		Title:       "Test Task",
		Description: "Task Description",
		AssigneeID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&task).Error
	assert.NoError(t, err)

	token, err := jwt.GenerateUserJwtToken(user.ID)
	assert.NoError(t, err)

	comment := map[string]interface{}{
		"comment": "This is a test comment",
		"task_id": task.ID,
	}
	jsonBody, err := json.Marshal(comment)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var responseComment models.Comment
	err = json.Unmarshal(resp.Body.Bytes(), &responseComment)
	assert.NoError(t, err)
	assert.Equal(t, "This is a test comment", responseComment.Comment)
	assert.Equal(t, task.ID, responseComment.TaskID)
	assert.Equal(t, user.ID, responseComment.UserID)
}

func TestGetCommentsByTaskID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	user := models.User{
		Username:     "commentuser",
		PasswordHash: "password",
	}
	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	task := models.Task{
		Title:       "Test Task",
		Description: "Task Description",
		AssigneeID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&task).Error
	assert.NoError(t, err)

	comments := []models.Comment{
		{
			Comment: "First comment",
			TaskID:  task.ID,
			UserID:  user.ID,
		},
		{
			Comment: "Second comment",
			TaskID:  task.ID,
			UserID:  user.ID,
		},
	}
	for _, comment := range comments {
		err = testutils.SharedTestDB.Create(&comment).Error
		assert.NoError(t, err)
	}

	token, err := jwt.GenerateUserJwtToken(user.ID)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/tasks/"+strconv.Itoa(int(task.ID))+"/comments", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseComments []models.Comment
	err = json.Unmarshal(resp.Body.Bytes(), &responseComments)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(responseComments))
}

func TestUpdateComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	user := models.User{
		Username:     "commentuser",
		PasswordHash: "password",
	}
	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	task := models.Task{
		Title:       "Test Task",
		Description: "Task Description",
		AssigneeID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&task).Error
	assert.NoError(t, err)

	comment := models.Comment{
		Comment: "Original comment",
		TaskID:  task.ID,
		UserID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&comment).Error
	assert.NoError(t, err)

	token, err := jwt.GenerateUserJwtToken(user.ID)
	assert.NoError(t, err)

	updateData := map[string]interface{}{
		"comment": "Updated comment",
	}
	jsonBody, err := json.Marshal(updateData)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/comments/"+strconv.Itoa(int(comment.ID)), bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseComment models.Comment
	err = json.Unmarshal(resp.Body.Bytes(), &responseComment)
	assert.NoError(t, err)
	assert.Equal(t, "Updated comment", responseComment.Comment)
}

func TestDeleteComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	// Create a user
	user := models.User{
		Username:     "commentuser",
		PasswordHash: "password",
	}
	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	task := models.Task{
		Title:       "Test Task",
		Description: "Task Description",
		AssigneeID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&task).Error
	assert.NoError(t, err)

	comment := models.Comment{
		Comment: "Test comment to delete",
		TaskID:  task.ID,
		UserID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&comment).Error
	assert.NoError(t, err)

	token, err := jwt.GenerateUserJwtToken(user.ID)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, "/api/v1/comments/"+strconv.Itoa(int(comment.ID)), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)

	var count int64
	testutils.SharedTestDB.Model(&models.Comment{}).Where("id = ?", comment.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteTask_CascadeComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	user := models.User{
		Username:     "taskdeleter",
		PasswordHash: "password",
	}
	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	task := models.Task{
		Title:       "Task To Delete",
		Description: "This task will be deleted",
		AssigneeID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&task).Error
	assert.NoError(t, err)

	for i := range 3 {
		comment := models.Comment{
			Comment: "Comment " + strconv.Itoa(i),
			TaskID:  task.ID,
			UserID:  user.ID,
		}
		err = testutils.SharedTestDB.Create(&comment).Error
		assert.NoError(t, err)
	}

	var commentCountBefore int64
	testutils.SharedTestDB.Model(&models.Comment{}).Where("task_id = ?", task.ID).Count(&commentCountBefore)
	assert.Equal(t, int64(3), commentCountBefore)

	token, err := jwt.GenerateUserJwtToken(user.ID)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, "/api/v1/tasks/"+strconv.Itoa(int(task.ID)), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Call GetTask to verify task was soft deleted
	req, err = http.NewRequest(http.MethodGet, "/api/v1/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	assert.NoError(t, err)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var commentCountAfter int64
	testutils.SharedTestDB.Model(&models.Comment{}).
		Joins("JOIN tasks ON tasks.id = comments.task_id").
		Where("tasks.deleted_at IS NULL").
		Count(&commentCountAfter)
	assert.Equal(t, int64(0), commentCountAfter)
}

func TestCommentRemaining_AfterUserDeletion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	user := models.User{
		Username:     "usertodelete",
		PasswordHash: "password",
	}
	err := testutils.SharedTestDB.Create(&user).Error
	assert.NoError(t, err)

	task := models.Task{
		Title:       "Task For Comment",
		Description: "Task description",
		AssigneeID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&task).Error
	assert.NoError(t, err)

	comment := models.Comment{
		Comment: "Comment from user to be deleted",
		TaskID:  task.ID,
		UserID:  user.ID,
	}
	err = testutils.SharedTestDB.Create(&comment).Error
	assert.NoError(t, err)

	err = testutils.SharedTestDB.Delete(&user).Error
	assert.NoError(t, err)

	var foundComment models.Comment
	err = testutils.SharedTestDB.First(&foundComment, comment.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, comment.ID, foundComment.ID)
	assert.Equal(t, comment.Comment, foundComment.Comment)
}
