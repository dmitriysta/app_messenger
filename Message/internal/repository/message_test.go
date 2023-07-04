//go:build integration
// +build integration

package repository

import (
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	entities "internal/entities"
)

var testMessage = &entities.Message{
	UserID:    1,
	ChannelID: 1,
	Content:   "Test message",
	CreatedAt: time.Now(),
}

func TestIntegration(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=testuser dbname=testdb sslmode=disable password=testpassword")
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	repo := &messageRepository{db: db, redisClient: redisClient}

	t.Run("CreateMessage", func(t *testing.T) {
		err := repo.CreateMessage(testMessage)
		assert.NoError(t, err)
	})

	t.Run("SendMessage", func(t *testing.T) {
		testMessage.Content = "New test message"
		err := repo.SendMessage(testMessage)
		assert.NoError(t, err)
	})

	t.Run("ReceiveMessage", func(t *testing.T) {
		message, err := repo.ReceiveMessage(testMessage.Id)
		assert.NoError(t, err)
		assert.Equal(t, testMessage.Content, message.Content)
	})

	t.Run("DeleteMessage", func(t *testing.T) {
		err := repo.DeleteMessage(testMessage.Id)
		assert.NoError(t, err)

		message, err := repo.ReceiveMessage(testMessage.Id)
		assert.Error(t, err)
		assert.Nil(t, message)
	})

	t.Run("GetMessages", func(t *testing.T) {
		messages, err := repo.GetMessages(testMessage.ChannelID, 0, 10)
		assert.NoError(t, err)
		assert.Greater(t, len(messages), 0)
	})
}
