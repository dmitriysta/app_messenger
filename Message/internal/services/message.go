package services

import (
	entities "internal/entities"
	repository "internal/repository"
)

type MessageService interface {
	CreateMessage(message *entities.Message) error
	SendMessage(message *entities.Message) error
	ReceiveMessage(id int) (*entities.Message, error)
	DeleteMessage(id int) error
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func (s *messageService) CreateMessage(message *entities.Message) error {
	if err := s.messageRepo.CreateMessage(message); err != nil {
		return err
	}
	return nil
}

func (s *messageService) SendMessage(message *entities.Message) error {
	if err := s.messageRepo.SendMessage(message); err != nil {
		return err
	}
	return nil
}

func (s *messageService) ReceiveMessage(id int) (*entities.Message, error) {
	if message, err := s.messageRepo.ReceiveMessage(id); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *messageService) DeleteMessage(id int) error {
	if err := s.messageRepo.DeleteMessage(id); err != nil {
		return err
	}
	return nil
}
