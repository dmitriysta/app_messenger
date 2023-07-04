//go:build integration
// +build integration

package repository

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	entities "internal/entities"
)

var testUser = &entities.User{
	Name:     "Test User",
	Email:    "test@email.com",
	Password: "password",
}

func TestIntegration(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb sslmode=disable password=testpassword")
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	repo := &userRepository{db: db}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(testUser)
		assert.NoError(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		user, err := repo.GetByID(testUser.Id)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Name, user.Name)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("Update", func(t *testing.T) {
		testUser.Name = "New Test User"
		err := repo.Update(testUser)
		assert.NoError(t, err)

		user, err := repo.GetByID(testUser.Id)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Name, user.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(testUser.Id)
		assert.NoError(t, err)

		user, err := repo.GetByID(testUser.Id)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
