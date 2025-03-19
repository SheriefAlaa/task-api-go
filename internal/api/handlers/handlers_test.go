package handlers_test

import (
	"os"
	"task-api-go/internal/testutils"
	"testing"
)

func TestMain(m *testing.M) {
	testutils.SetupTestDB()
	code := m.Run()
	testutils.CleanupTestDB()

	os.Exit(code)
}
