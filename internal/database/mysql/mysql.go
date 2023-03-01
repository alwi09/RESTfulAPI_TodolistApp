package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"todolist_gin_gorm/internal/config"
	"todolist_gin_gorm/internal/model/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	
	fmt.Printf("%+v\n ", cfg)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel: logger.Info,
				Colorful: true,
			},
		),
	})

	if err != nil {
		panic("Cannot connect to database")
	}

	err = db.AutoMigrate(&entity.Todos{})
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Connect database successfully")

	return db,err
}