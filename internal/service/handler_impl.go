package service

import (
	"errors"
	"net/http"
	"strconv"
	"todolist_gin_gorm/internal/database"
	"todolist_gin_gorm/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandlerImpl struct {
	todolistRepository *database.TodoRepository
}

func NewHandlerImpl(repository *database.TodoRepository) *HandlerImpl {
	return &HandlerImpl{
		todolistRepository: repository,
	}
}

func (handler *HandlerImpl) CreateHandlerTodolist(ctx *gin.Context) {
	todos := new(dto.CreateTodolistRequest)
	err := ctx.ShouldBindJSON(todos)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	newList, errCreate := handler.todolistRepository.Create(todos.Title, todos.Description)
	if errCreate != nil {
		logrus.Error(errCreate.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: errCreate.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	logrus.Info(http.StatusCreated, "create todolist successfully", todos)
	ctx.JSON(http.StatusCreated, dto.TodolistResponseCreate{

		Message: "create todolist successfully",
		Status:  http.StatusCreated,
		Data:    newList,
	})

	return
}

func (handler *HandlerImpl) GetAllHandlerTodolist(ctx *gin.Context) {
	todos, err := handler.todolistRepository.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}

	logrus.Info(http.StatusOK, "get all todolist successfully", todos)
	ctx.JSON(http.StatusOK, dto.TodolistResponseGetAll{

		Message: "get all todolist successfully",
		Status: http.StatusOK,
		More: len(todos),
		Data: todos,
	})

	return
}

func (handler *HandlerImpl) GetIDHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}

	todos, err := handler.todolistRepository.GetID(todoID)
	if err != nil {
		logrus.Errorf("failed whe get todolist bi id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}

	if todos == nil {
		logrus.Error(http.StatusNotFound, errors.New("todolist by id not found"))
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{

			Message: "todolist by id not found",
			Status: http.StatusNotFound,
		})
		return
	}

	logrus.Info(http.StatusOK, "get todolist by id successfully")
	ctx.AbortWithStatusJSON(http.StatusOK, dto.TodolistResponseGetID{
		
		Message: "get todolist by id successfully",
		Status: http.StatusOK,
		Data: *todos,
	})

	return

}

func (handler *HandlerImpl) UpdateHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}

	todos := new(dto.UpdateTodolistRequest)
	err = ctx.ShouldBindJSON(todos)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}

	id, err := handler.todolistRepository.GetID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}

	if id == nil {
		logrus.Error(http.StatusNotFound, errors.New("todolist by id not not found"))
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{

			Message: "todolist by id not found",
			Status: http.StatusNotFound,
		})
		return
	}

	update, err := handler.todolistRepository.Update(todoID, todos.RequestUpdateTodolist())
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}

	if update == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, dto.TodolistResponseGetID{

			Message: "not change",
			Status: http.StatusOK,
		})
		return
	}

	logrus.Info(http.StatusOK, "update todolist successfully")
	ctx.JSON(http.StatusOK, dto.TodolistResponseUpdate{

		Message: "update todolist successfully",
		Status: http.StatusOK,
		Data: todos,
	})
	return
}

func (handler *HandlerImpl) DeleteHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}

	IDtodo, err := handler.todolistRepository.Delete(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{

			Message: err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}

	if IDtodo == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{

			Message: "id not found",
			Status: http.StatusNotFound,
		})
		return
	}

	logrus.Info(http.StatusOK, "delete todolist successfully")
	ctx.JSON(http.StatusOK, dto.TodolistResponseDelete{

		Message: "delete todolist successfully",
		Status: http.StatusOK,
	})

	return
}