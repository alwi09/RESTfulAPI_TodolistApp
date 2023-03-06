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
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": err.Error(),
				"status":  http.StatusBadRequest,
			})
		return
	}

	if _, errCreate := handler.todolistRepository.Create(todos.Title, todos.Description); errCreate != nil {
		logrus.Error(errCreate.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": errCreate.Error(),
				"status":  http.StatusBadRequest,
			})
		return
	}

	logrus.Info(http.StatusCreated, "create todolist successfully", todos)
	ctx.JSON(http.StatusCreated,
		gin.H{
			"message": "create todolist successfully",
			"status":  http.StatusCreated,
			"data":    todos,
		})

	return
}

func (handler *HandlerImpl) GetAllHandlerTodolist(ctx *gin.Context) {
	todos, err := handler.todolistRepository.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"message": err.Error(),
				"status":  http.StatusInternalServerError,
			})
		return
	}

	logrus.Info(http.StatusOK, "Get all successfully", todos)
	ctx.JSON(http.StatusOK,
		gin.H{
			"message": "Get all successfully",
			"status":  http.StatusOK,
			"more":    len(todos),
			"data":    todos,
		})

	return
}

func (handler *HandlerImpl) GetIDHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": err,
				"status":  http.StatusBadRequest,
			})
		return
	}

	todos, err := handler.todolistRepository.GetID(todoID)
	if err != nil {
		logrus.Errorf("failed whe get todolist bi id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"message": err,
				"status":  http.StatusInternalServerError,
			})
		return
	}

	if todos == nil {
		logrus.Error(http.StatusNotFound, errors.New("todolist by id not found"))
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"message": "todolist by id not found",
				"status":  http.StatusNotFound,
			})
		return
	}

	logrus.Info(http.StatusOK, "get todolist by id successfully")
	ctx.AbortWithStatusJSON(http.StatusOK,
		gin.H{
			"message": "get todolist by id successfully",
			"status":  http.StatusOK,
			"data":    todos,
		})

	return

}

func (handler *HandlerImpl) UpdateHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": err,
				"status":  http.StatusBadRequest,
			})
		return
	}

	todos := new(dto.UpdateTodolistRequest)
	err = ctx.ShouldBindJSON(todos)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": err,
				"status":  http.StatusBadRequest,
			})
		return
	}

	id, err := handler.todolistRepository.GetID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"message": err,
				"status":  http.StatusInternalServerError,
			})
		return
	}

	if id == nil {
		logrus.Error(http.StatusNotFound, errors.New("todolist by id not not found"))
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"message": "todolist by id not found",
				"status":  http.StatusNotFound,
			})
		return
	}

	update, err := handler.todolistRepository.Update(todoID, todos.RequestUpdateTodolist())
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"message": err,
				"status":  http.StatusInternalServerError,
			})
		return
	}

	if update == 0 {
		ctx.AbortWithStatusJSON(http.StatusOK,
			gin.H{
				"message": err,
				"status":  http.StatusOK,
			})
		return
	}

	logrus.Info(http.StatusOK, "update todolist successfully")
	ctx.JSON(http.StatusOK,
		gin.H{
			"message": "update data successfully",
			"status":  http.StatusOK,
			"data":    todos,
		})
	return
}

func (handler *HandlerImpl) DeleteHandlerTodolist(ctx *gin.Context) {
	todolistId := ctx.Param("todolistId")
	todoID, err := strconv.ParseInt(todolistId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": err,
				"status":  http.StatusBadRequest,
			})
		return
	}

	IDtodo, err := handler.todolistRepository.Delete(todoID)
	if err != nil {
		logrus.Errorf("failed when get todolist by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"message": err,
				"status":  http.StatusInternalServerError,
			})
		return
	}

	if IDtodo == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"message": err,
				"status":  http.StatusNotFound,
			})
		return
	}

	logrus.Info(http.StatusOK, "delete todolist successfully")
	ctx.JSON(http.StatusOK,
		gin.H{
			"message": "delete todolist successfully",
			"status":  http.StatusOK,
		})

	return
}
