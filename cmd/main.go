package main

import (
	"context"
	"errors"
	"io"

	"os"
	"todolist_gin_gorm/cmd/router"
	"todolist_gin_gorm/internal/config"
	"todolist_gin_gorm/internal/database"
	"todolist_gin_gorm/internal/database/mysql"
	"todolist_gin_gorm/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func SetupLogOutput() {
	file, _ := os.Create("gin-log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}

func main() {

	SetupLogOutput()

	ctx := context.Background()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	var cfg config.Config

	// initialize config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logrus.Fatal(errors.New("Error"), err)
	}

	// initialize db connect
	db, err := mysql.Connect(ctx, &cfg)
	if err != nil {
		return
	}

	// initialize migration
	err = mysql.Migrate(db)
	if err != nil {
		logrus.Fatalf("error running schema migration %v", err)
	}

	// initialize repositories
	todolistRepository := database.NewTodoRepository(db)

	// initialize service
	todolistHandler := service.NewHandlerImpl(todolistRepository)

	// initialize router
	routeBuilder := router.NewRouteBuilder(todolistHandler)
	routerInit := routeBuilder.RouteInit()

	err = routerInit.Run(":1234")
	if err != nil {
		logrus.Fatal(err)
	}
}
