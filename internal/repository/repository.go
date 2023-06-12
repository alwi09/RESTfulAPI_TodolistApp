package repository

import "todolist_gin_gorm/internal/model/entity"

type Repository interface {
	GetAll() ([]entity.Todos, error)
	GetID(todoID int64) (*entity.Todos, error)
	Create(title string, description string) (*entity.Todos, error)
	Update(todoID int64, updates map[string]interface{}) (*entity.Todos, error)
	Delete(todoID int64) (int64, error)
	CreateUser(user *entity.Users) error
	FindUserByEmail(username string) (*entity.Users, error)
}
