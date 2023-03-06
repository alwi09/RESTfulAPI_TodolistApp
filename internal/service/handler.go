package service

import "github.com/gin-gonic/gin"

type Handler interface {
	CreateHandlerTodolist(ctx *gin.Context)
	GetAllHandlerTodolist(ctx *gin.Context)
	GetIDHandlerTodolist(ctx *gin.Context)
	UpdateHandlerTodolist(ctx *gin.Context)
	DeleteHandlerTodolist(ctx *gin.Context)
}