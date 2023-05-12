package repository

import (
	"encoding/json"
	"fmt"

	entities "internal/entities"
)

type RedisRepository struct {
	Client *redis.Client
}

func (r *RedisRepository) GetMessages(channelID string) ([]*entities.Message, error) {

	exists, err := r.Client.Exists(channelID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check if channelID exists in Redis: %v", err)
	}
	if exists == 0 {
		return nil, nil
	}

	val, err := r.Client.Get(channelID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from Redis: %v", err)
	}

	var messages []*entities.Message
	err = json.Unmarshal([]byte(val), &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages from Redis: %v", err)
	}

	return messages, nil
}

func (r *RedisRepository) SetMessages(channelID string, messages []*entities.Message) error {
	val, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal messages to JSON: %v", err)
	}

	err = r.Client.Set(channelID, val, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set messages in Redis: %v", err)
	}

	return nil
}
