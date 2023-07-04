package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserController(t *testing.T) {
	userController := &UserController{}

	t.Run("GetUserByID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/user?id=1", nil)
		rec := httptest.NewRecorder()

		userController.GetUserByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res GetUserResponce
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)

	})

	t.Run("CreateUser", func(t *testing.T) {
		user := GetUserResponce{
			Name:     "Test",
			Email:    "test@example.com",
			Password: "password",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		userController.CreateUser(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user := GetUserResponce{
			Id:       1,
			Name:     "Test Updated",
			Email:    "test_updated@example.com",
			Password: "newpassword",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		userController.UpdateUser(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/user?id=1", nil)
		rec := httptest.NewRecorder()

		userController.DeleteUser(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
