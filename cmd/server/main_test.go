package main_test

import (
	"net/http"
	"net/http/httptest"
	"task-api-go/internal/testutils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	testutils.CascadeDB(t)

	r := testutils.SetupTestRouter(testutils.SharedTestDB)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	assert.NoError(t, err)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"status":"API is healthy!"}`
	assert.JSONEq(t, expected, rr.Body.String())

}
