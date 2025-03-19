// internal/testutils/shared_test_db.go
package testutils

import (
	"fmt"
	"log"
	"sync"
	"task-api-go/internal/models"
	"testing"

	"github.com/ory/dockertest/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	SharedTestDB *gorm.DB
	pool         *dockertest.Pool
	resource     *dockertest.Resource
	once         sync.Once
	initialized  bool
)

func SetupTestDB() {
	once.Do(func() {
		var err error
		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}

		resource, err = pool.Run("postgres", "16-alpine", []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=testdb",
		})
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}

		if err = pool.Retry(func() error {
			var err error
			dsn := fmt.Sprintf("host=localhost port=%s user=postgres password=postgres dbname=testdb sslmode=disable",
				resource.GetPort("5432/tcp"))
			SharedTestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return err
			}
			return SharedTestDB.Exec("SELECT 1").Error
		}); err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}

		SharedTestDB.AutoMigrate(&models.User{}, &models.Task{})

		initialized = true
		log.Println("Shared test database initialized")
	})
}

func CascadeDB(t *testing.T) {
	t.Helper()
	if !initialized {
		SetupTestDB()
	}

	tables := []string{"users", "tasks"}
	for _, table := range tables {
		SharedTestDB.Exec("TRUNCATE TABLE " + table + " CASCADE")
	}
}

func CleanupTestDB() {
	if resource != nil && pool != nil {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}
}
