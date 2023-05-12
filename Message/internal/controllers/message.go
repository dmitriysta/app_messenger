package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	messageRepo repository.MessageRepository
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
