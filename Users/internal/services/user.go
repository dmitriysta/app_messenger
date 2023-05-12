package services

import (
	entities "internal/entities"
	repository "internal/repository"
)

type UserService interface {
	CreateUser(user *entities.User) error
	UpdateUser(user *entities.User) error
	DeleteUser(id int) error
	GetUserByID(id int) (*entities.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *entities.User) error {
	err := s.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(user *entities.User) error {
	err := s.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUser(id int) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUserByID(id int) (*entities.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
