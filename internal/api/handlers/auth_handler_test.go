package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-api-go/internal/api/requests"
	"task-api-go/internal/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignupUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	reqBody := requests.SignupInput{
		Username:             "testuser",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}
	jsonBody, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]any
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "username")
	assert.Equal(t, "testuser", response["username"])
	assert.Contains(t, response, "ID")
}

func TestSignupUser_Failure_WhenPasswordDoesNotMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	reqBody := requests.SignupInput{
		Username:             "testuser",
		Password:             "password123",
		PasswordConfirmation: "passwordzzz",
	}
	jsonBody, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]any
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Passwords do not match", response["error"])
}
