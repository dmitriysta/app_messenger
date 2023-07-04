package services

import (
	"github.com/golang/mock/gomock"
	"internal/entities"
	"internal/repository/mocks"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)

	msg := &entities.Message{Id: 1, Content: "Hello, Test!"}

	mockRepo.EXPECT().CreateMessage(msg).Return(nil)

	service := &messageService{
		messageRepo: mockRepo,
	}

	err := service.CreateMessage(msg)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)

	msg := &entities.Message{Id: 1, Content: "Hello, Test!"}

	mockRepo.EXPECT().SendMessage(msg).Return(nil)

	service := &messageService{
		messageRepo: mockRepo,
	}

	err := service.SendMessage(msg)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReceiveMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)

	msg := &entities.Message{Id: 1, Content: "Hello, Test!"}

	mockRepo.EXPECT().ReceiveMessage(1).Return(msg, nil)

	service := &messageService{
		messageRepo: mockRepo,
	}

	receivedMsg, err := service.ReceiveMessage(1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if receivedMsg != msg {
		t.Errorf("unexpected message: got %v, want %v", receivedMsg, msg)
	}
}

func TestDeleteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepository(ctrl)

	mockRepo.EXPECT().DeleteMessage(1).Return(nil)

	service := &messageService{
		messageRepo: mockRepo,
	}

	err := service.DeleteMessage(1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
