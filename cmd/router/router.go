package router

import (
	"todolist_gin_gorm/internal/middleware"
	"todolist_gin_gorm/internal/service"

	"github.com/gin-gonic/gin"
)

type RouteBuilder struct {
	todoHandler *service.HandlerImpl
}

func NewRouteBuilder(todoHandler *service.HandlerImpl) *RouteBuilder {
	return &RouteBuilder{
		todoHandler: todoHandler,
	}
}

func (routeBuilder *RouteBuilder) RouteInit() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery(), middleware.Logger(), middleware.XAPIKEY())

	router.GET("/api/find_all_todolist", routeBuilder.todoHandler.GetAllHandlerTodolist)
	router.GET("/api/find_by_id_todolist/:todolistId", routeBuilder.todoHandler.GetIDHandlerTodolist)
	router.POST("/api/create_todolist", routeBuilder.todoHandler.CreateHandlerTodolist)
	router.PUT("/api/update_todolist/:todolistId", routeBuilder.todoHandler.UpdateHandlerTodolist)
	router.DELETE("/api/delete_todolist/:todolistId", routeBuilder.todoHandler.DeleteHandlerTodolist)

	return router
}