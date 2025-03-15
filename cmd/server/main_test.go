package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	// Create a request to send to the endpoint
	req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"status":"API is healthy!"}`
	assert.JSONEq(t, expected, rr.Body.String())

}
