package database

import (
	"errors"
	"todolist_gin_gorm/internal/model/entity"

	"gorm.io/gorm"
)

//adaptop pattern

type TodoRepository struct {
	DB *gorm.DB
}

func NewTodoRepository(DB *gorm.DB) *TodoRepository {
	return &TodoRepository{
		DB: DB,
	}
}

func (repository *TodoRepository) Create(title string, description string) (*entity.Todos, error) {
	todos := entity.Todos{
		Title:       title,
		Description: description,
	}

	result := repository.DB.Create(&todos)
	return &todos, result.Error
}

func (repository *TodoRepository) GetAll() ([]entity.Todos, error) {
	var todos []entity.Todos
	result := repository.DB.Find(&todos)

	return todos, result.Error
}

func (repository *TodoRepository) GetID(todoID int64) (*entity.Todos, error) {
	var todos entity.Todos
	result := repository.DB.Where("todos_id = ?", todoID).First(&todos)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &todos, result.Error
}

func (repository *TodoRepository) Update(todoID int64, updates map[string]interface{}) (*entity.Todos, error) {
	var todos entity.Todos
	result := repository.DB.Model(&todos).Where("todos_id = ?", todoID).Updates(updates)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &todos, result.Error
}

func (repository *TodoRepository) Delete(todoID int64) (int64, error) {
	todos := entity.Todos{
		Id: todoID,
	}

	result := repository.DB.Delete(&todos)

	return result.RowsAffected, result.Error
}

func (repository *TodoRepository) CreateUser(user *entity.Users) error {
	if err := repository.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (repository *TodoRepository) FindUserByEmail(email string) (*entity.Users, error) {
	var user entity.Users
	if err := repository.DB.Where("email = ? ", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
