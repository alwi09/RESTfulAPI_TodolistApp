package dto

import (
	"errors"
	"todolist_gin_gorm/internal/model/entity"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=2"`
}

// untuk validate register request
func ValidateRegisterRequest(user *entity.Users) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// untuk validate login request
func ValidateLoginRequest(user *entity.Users) error {
	if len(user.Email) < 2 || len(user.Password) < 2 {
		return errors.New("email and password must be at least 2 characters long")
	}

	return nil
}
