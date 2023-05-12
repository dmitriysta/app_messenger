package repository

import (
	"database/sql"
	"fmt"
	entities "internal/entities"
	"log"
)

type MessageRepository interface {
	CreateMessage(message *entities.Message) error
	SendMessage(message *entities.Message) error
	ReceiveMessage(id int) (*entities.Message, error)
	DeleteMessage(id int) error
	GetMessages(channelID int, offset int, limit int) ([]*entities.Message, error)
}

type messageRepository struct {
	db          *sql.DB
	redisClient *redis.Client
}

func (r *messageRepository) CreateMessage(message *entities.Message) error {
	sqlQuery := "INSERT INTO messages(user_id, channel_id, content, created_at) VALUES($1, $2, $3, $4) RETURNING id"
	if err := r.db.QueryRow(sqlQuery, message.UserID, message.ChannelID, message.Content, message.CreatedAt).Scan(&message.Id); err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}
	return nil
}

func (r *messageRepository) SendMessage(message *entities.Message) error {
	sqlQuery := "UPDATE messages SET user_id = $1, channel_id = $2, content = $3, created_at = $4 WHERE id = $5"
	if _, err := r.db.Exec(sqlQuery, message.UserID, message.ChannelID, message.Content, message.CreatedAt, message.Id); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func (r *messageRepository) ReceiveMessage(id int) (*entities.Message, error) {
	sqlQuery := "SELECT id, user_id, channel_id, content, created_at FROM messages WHERE id = $1"
	message := &entities.Message{}

	if err := r.db.QueryRow(sqlQuery, id).Scan(&message.Id, &message.UserID, &message.ChannelID, &message.Content, &message.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to receive message: %v", err)
	}
	return message, nil
}

func (r *messageRepository) DeleteMessage(id int) error {
	sqlQuery := "DELETE FROM messages WHERE id = $1"
	if _, err := r.db.Exec(sqlQuery, id); err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}

func (r *messageRepository) GetMessages(channelID int, offset int, limit int) ([]*entities.Message, error) {
	redisKey := fmt.Sprintf("messages:%d:%d:%d", channelID, offset, limit)

	cachedMessages, err := r.redisClient.GetMessages(redisKey)
	if err == nil {
		return cachedMessages, nil
	}

	sqlQuery := "SELECT id, user_id, channel_id, content, created_at FROM messages WHERE channel_id = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3"
	rows, err := r.db.Query(sqlQuery, channelID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}
	defer rows.Close()

	messages := []*entities.Message{}
	for rows.Next() {
		message := &entities.Message{}
		if err := rows.Scan(&message.Id, &message.UserID, &message.ChannelID, &message.Content, &message.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %v", err)
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over messages rows: %v", err)
	}

	if err := r.redisClient.SetMessages(redisKey, messages); err != nil {
		log.Printf("Failed to set messages to Redis cache: %v\n", err)
	}

	return messages, nil
}
