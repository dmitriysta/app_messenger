package services

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"internal/entities"
	"internal/repository/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	testUser := &entities.User{ID: 1, Name: "Test User"}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().Create(testUser).Return(nil)

		service := NewUserService(mockUserRepo)
		err := service.CreateUser(testUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockUserRepo.EXPECT().Create(testUser).Return(errors.New("Panic error"))

		service := NewUserService(mockUserRepo)
		err := service.CreateUser(testUser)

		assert.Error(t, err)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	testUser := &entities.User{ID: 1, Name: "Test User"}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().Update(testUser).Return(nil)

		service := NewUserService(mockUserRepo)
		err := service.UpdateUser(testUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockUserRepo.EXPECT().Update(testUser).Return(errors.New("Panic error"))

		service := NewUserService(mockUserRepo)
		err := service.UpdateUser(testUser)

		assert.Error(t, err)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	testUserID := 1

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().Delete(testUserID).Return(nil)

		service := NewUserService(mockUserRepo)
		err := service.DeleteUser(testUserID)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockUserRepo.EXPECT().Delete(testUserID).Return(errors.New("Panic error"))

		service := NewUserService(mockUserRepo)
		err := service.DeleteUser(testUserID)

		assert.Error(t, err)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	testUserID := 1
	testUser := &entities.User{ID: testUserID, Name: "Test User"}

	t.Run("успех", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(testUserID).Return(testUser, nil)

		service := NewUserService(mockUserRepo)
		user, err := service.GetUserByID(testUserID)

		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
	})

	t.Run("ошибка", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(testUserID).Return(nil, errors.New("принудительная ошибка"))

		service := NewUserService(mockUserRepo)
		user, err := service.GetUserByID(testUserID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
