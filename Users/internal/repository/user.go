package repository

import (
	"database/sql"
	"fmt"
	entities "internal/entities"
)

type UserRepository interface {
	Create(user *entities.User) error
	Update(user *entities.User) error
	Delete(id int) error
	GetByID(id int) (*entities.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) Create(user *entities.User) error {
	sqlQuery := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(sqlQuery, user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (r *userRepository) GetByID(id int) (*entities.User, error) {
	sqlQuery := "SELECT from users (id, name, email, password) WHERE id = $1"
	user := &entities.User{}
	err := r.db.QueryRow(sqlQuery, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

func (r *userRepository) Update(user *entities.User) error {
	sqlQuery := "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4"
	err := r.db.QueryRow(sqlQuery, user.Name, user.Email, user.Password, user.Id)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (r *userRepository) Delete(id int) error {
	sqlQuery := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(sqlQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
