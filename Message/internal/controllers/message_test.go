package controllers

import (
	"bytes"
	"encoding/json"
	"internal/repository/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	mockRepo := mocks.NewMockMessageRepository()

	controller := MessageController{messageRepo: mockRepo}

	message := GetMessageResponce{
		Id:        1,
		UserID:    1,
		ChannelID: 1,
		Content:   "Hello, World!",
	}

	mockRepo.EXPECT().CreateMessage(gomock.Any()).Return(nil)

	body, _ := json.Marshal(message)
	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	controller.CreateMessage(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	var response GetMessageResponce
	json.NewDecoder(recorder.Body).Decode(&response)

	if response != message {
		t.Errorf("handler returned unexpected body: got %v, want %v", response, message)
	}
}

func TestSendMessage(t *testing.T) {
	mockRepo := mocks.NewMockMessageRepository()

	controller := MessageController{messageRepo: mockRepo}

	message := GetMessageResponce{
		Id:        1,
		UserID:    1,
		ChannelID: 1,
		Content:   "Hello, World!",
	}

	mockRepo.EXPECT().SendMessage(gomock.Any()).Return(nil)

	body, _ := json.Marshal(message)
	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	controller.SendMessage(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	var response GetMessageResponce
	json.NewDecoder(recorder.Body).Decode(&response)

	if response != message {
		t.Errorf("handler returned unexpected body: got %v, want %v", response, message)
	}
}

func TestReceiveMessage(t *testing.T) {
	mockRepo := mocks.NewMockMessageRepository()

	controller := MessageController{messageRepo: mockRepo}

	message := GetMessageResponce{
		Id:        1,
		UserID:    1,
		ChannelID: 1,
		Content:   "Hello, World!",
	}

	mockRepo.EXPECT().ReceiveMessage(1).Return(&message, nil)

	req, _ := http.NewRequest("GET", "/message?id=1", nil)
	recorder := httptest.NewRecorder()

	controller.ReceiveMessage(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	var response GetMessageResponce
	json.NewDecoder(recorder.Body).Decode(&response)

	if response != message {
		t.Errorf("handler returned unexpected body: got %v, want %v", response, message)
	}
}

func TestDeleteMessage(t *testing.T) {
	mockRepo := mocks.NewMockMessageRepository()

	controller := MessageController{messageRepo: mockRepo}

	id := 1

	mockRepo.EXPECT().DeleteMessage(id).Return(nil)

	req, _ := http.NewRequest("DELETE", "/message?id=1", nil)
	recorder := httptest.NewRecorder()

	controller.DeleteMessage(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}
}
