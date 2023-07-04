package controllers

import (
	"internal/entities"
	"internal/repository/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userController := UserController{userRepo: mockUserRepo}

	user := &entities.User{Id: 1, Name: "Test User", Email: "test@email.com", Password: "password"}

	mockUserRepo.EXPECT().GetByID(1).Return(user, nil)

	req, err := http.NewRequest("GET", "/user?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.GetUserByID)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"Test User","email":"test@email.com","password":"password"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userController := UserController{userRepo: mockUserRepo}

	user := &entities.User{Id: 1, Name: "Test User", Email: "test@email.com", Password: "password"}

	mockUserRepo.EXPECT().Create(user).Return(nil)

	userJson := `{"id":1,"name":"Test User","email":"test@email.com","password":"password"}`
	req, err := http.NewRequest("POST", "/user", strings.NewReader(userJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.CreateUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestUpdateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userController := UserController{userRepo: mockUserRepo}

	user := &entities.User{Id: 1, Name: "Test User", Email: "test@email.com", Password: "password"}

	mockUserRepo.EXPECT().Update(user).Return(nil)

	userJson := `{"id":1,"name":"Test User","email":"test@email.com","password":"password"}`
	req, err := http.NewRequest("PUT", "/user", strings.NewReader(userJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.UpdateUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userController := UserController{userRepo: mockUserRepo}

	userId := 1

	// Указываем, что мы ожидаем вызова метода Delete с определенным аргументом и что он должен вернуть.
	mockUserRepo.EXPECT().Delete(userId).Return(nil)

	req, err := http.NewRequest("DELETE", "/user?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userController.DeleteUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
