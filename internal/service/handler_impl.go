package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todolist_gin_gorm/internal/config"
	"todolist_gin_gorm/internal/model/dto"
	"todolist_gin_gorm/internal/model/entity"
	"todolist_gin_gorm/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type HandlerImpl struct {
	todolistRepository repository.Repository
}

func NewHandlerImpl(repository repository.Repository) *HandlerImpl {
	return &HandlerImpl{
		todolistRepository: repository,
	}
}

// RegisterHandler handles the registration request
func (handler *HandlerImpl) RegisterHandler(ctx *gin.Context) {
	// Parse request body
	var user entity.Users
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "failed to create user",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// validate request register
	err := dto.ValidateRegisterRequest(&user)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Check if user already exists
	existingUser, err := handler.todolistRepository.FindUserByEmail(user.Email)
	if existingUser != nil {
		logrus.Warn("user already exist", err.Error())
		ctx.AbortWithStatusJSON(http.StatusConflict, dto.ErrorResponse{
			Message: "user already exists",
			Status:  http.StatusConflict,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "an error occurred",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Create user
	newUser := &entity.Users{
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
	}
	err = handler.todolistRepository.CreateUser(newUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "failed to create user",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Return success message
	ctx.JSON(http.StatusOK, dto.CreateUserResponse{
		Message: "user created successfully",
		Status:  http.StatusOK,
	})
}

// LoginHandler handles the login request
func (handler *HandlerImpl) LoginHandler(ctx *gin.Context) {
	// Parse request body
	var user entity.Users
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	// validate request login
	err := dto.ValidateLoginRequest(&user)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Find user by Email
	existingUser, err := handler.todolistRepository.FindUserByEmail(user.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "an error occurred",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if existingUser == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid credentials",
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid credentials",
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Create JWT token
	tokenString, err := config.CreateJWTToken(user.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create token",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Return token in response
	ctx.JSON(http.StatusOK, dto.UserLoginResponse{
		Message: fmt.Sprintf("hello %s! you are now logged in", user.Username),
		Status:  http.StatusOK,
		Token:   tokenString,
	})
}

func (handler *HandlerImpl) CreateHandlerTodolist(ctx *gin.Context) {
	todos := new(dto.CreateTodolistRequest)
	err := ctx.ShouldBindJSON(todos)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: "invalid input validation",
			Status:  http.StatusBadRequest,
		})
		return
	}

	newList, errCreate := handler.todolistRepository.Create(todos.Title, todos.Description)
	if errCreate != nil {
		logrus.Error(errCreate.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	logrus.Info(http.StatusCreated, "create todolist successfully", todos)
	ctx.JSON(http.StatusCreated, dto.TodolistResponseCreate{

		Message: "create todolist successfully",
		Status:  http.StatusCreated,
		Data:    *newList,
	})

}

func (handler *HandlerImpl) GetAllHandlerTodolist(ctx *gin.Context) {
	todos, err := handler.todolistRepository.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	logrus.Info(http.StatusOK, "get all todolist successfully", todos)
	ctx.JSON(http.StatusOK, dto.TodolistResponseGetAll{

		Message: "get all todolist successfully",
		Status:  http.StatusOK,
		More:    len(todos),
		Data:    todos,
	})

}

func (handler *HandlerImpl) GetIDHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	todos, err := handler.todolistRepository.GetID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	if todos == nil {
		logrus.Error(http.StatusNotFound, errors.New("todolist by id not found"))
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{

			Message: "todolist by id not found",
			Status:  http.StatusNotFound,
		})
		return
	}

	logrus.Info(http.StatusOK, "get todolist by id successfully")
	ctx.AbortWithStatusJSON(http.StatusOK, dto.TodolistResponseGetID{

		Message: "get todolist by id successfully",
		Status:  http.StatusOK,
		Data:    *todos,
	})

}

func (handler *HandlerImpl) UpdateHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	todos := new(dto.UpdateTodolistRequest)
	err = ctx.ShouldBindJSON(todos)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	id, err := handler.todolistRepository.GetID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	if id == nil {
		logrus.Error(http.StatusNotFound, errors.New("todolist by id not not found"))
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{

			Message: "todolist by id not found",
			Status:  http.StatusNotFound,
		})
		return
	}

	update, err := handler.todolistRepository.Update(todoID, todos.RequestUpdateTodolist())
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	if update == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, dto.TodolistResponseGetID{

			Message: "not change",
			Status:  http.StatusOK,
		})
		return
	}

	logrus.Info(http.StatusOK, "update todolist successfully")
	ctx.JSON(http.StatusOK, dto.TodolistResponseUpdate{

		Message: "update todolist successfully",
		Status:  http.StatusOK,
		Data:    todos,
	})

}

func (handler *HandlerImpl) DeleteHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	IDtodo, err := handler.todolistRepository.Delete(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	if IDtodo == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{

			Message: "id not found",
			Status:  http.StatusNotFound,
		})
		return
	}

	logrus.Info(http.StatusOK, "delete todolist successfully")
	ctx.JSON(http.StatusOK, dto.TodolistResponseDelete{

		Message: "delete todolist successfully",
		Status:  http.StatusOK,
	})

}
