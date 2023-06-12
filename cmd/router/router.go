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
	// create Gin - Router
	router := gin.Default()

	// router.Use(gin.Recovery(), middleware.Logger(), middleware.XAPIKEY())
	router.Use(gin.Recovery(), middleware.Logger())

	// Group routes that require authentication
	authGroup := router.Group("/api")
	authGroup.Use(middleware.AuthMiddlewareJWT()) // Apply authentication middleware

	// register routes
	authGroup.GET("/find_all_todolist", routeBuilder.todoHandler.GetAllHandlerTodolist)
	authGroup.GET("/find_by_id_todolist/:todolistId", routeBuilder.todoHandler.GetIDHandlerTodolist)
	authGroup.POST("/create_todolist", routeBuilder.todoHandler.CreateHandlerTodolist)
	authGroup.PUT("/update_todolist/:todolistId", routeBuilder.todoHandler.UpdateHandlerTodolist)
	authGroup.DELETE("/delete_todolist/:todolistId", routeBuilder.todoHandler.DeleteHandlerTodolist)

	// public routes
	router.POST("/register", routeBuilder.todoHandler.RegisterHandler)
	router.POST("/login", routeBuilder.todoHandler.LoginHandler)

	return router
}
