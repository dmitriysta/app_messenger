package controllers

import (
	"encoding/json"
	"messenger/Message/internal/repository/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	repository "internal/repository"
)

type GetMessageResponce struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ChannelID int       `json:"channel_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageController struct {
	messageRepo   repository.MessageRepository
	kafkaProducer *kafka.Writer
}

func NewMessageController(logger *zap.Logger, producer *kafka.Writer) *MessageController {
	return &MessageController{
		messageRepo:   repository.NewMessageRepository(),
		kafkaProducer: producer,
	}
}

func (c *MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message GetMessageResponce
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message.CreatedAt = time.Now()

	if err := c.messageRepo.CreateMessage(&message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = c.kafkaProducer.WriteMessages(
		kafka.Message{
			Value: messageBytes,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (c *MessageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message GetMessageResponce
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.messageRepo.SendMessage(&message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func (c *MessageController) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := c.messageRepo.ReceiveMessage(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(message)
}

func (c *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.messageRepo.DeleteMessage(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := c.messageRepo.GetMessages(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func TestGetMessages(t *testing.T) {
	mockRepo := mocks.NewMockMessageRepository()

	controller := MessageController{messageRepo: mockRepo}

	page := 1
	limit := 10

	messages := []GetMessageResponce{
		{Id: 1, UserID: 1, ChannelID: 1, Content: "Hello, World!"},
		{Id: 2, UserID: 2, ChannelID: 1, Content: "Welcome!"},
	}

	mockRepo.EXPECT().GetMessages(page, limit).Return(messages, nil)

	req, _ := http.NewRequest("GET", "/messages?page=1&limit=10", nil)
	recorder := httptest.NewRecorder()

	controller.GetMessages(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	var response []GetMessageResponce
	json.NewDecoder(recorder.Body).Decode(&response)

	if len(response) != len(messages) {
		t.Errorf("handler returned wrong number of messages: got %v, want %v", len(response), len(messages))
	}

	for i, msg := range messages {
		if response[i] != msg {
			t.Errorf("handler returned unexpected message: got %v, want %v", response[i], msg)
		}
	}
}
