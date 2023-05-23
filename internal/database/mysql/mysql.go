package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"todolist_gin_gorm/internal/config"

	"github.com/golang-migrate/migrate/v4"
	mysqlMigration "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})

	if err != nil {
		panic("Cannot connect to database")
	}

	// err = db.AutoMigrate(&entity.Todos{})
	// if err != nil {
	// 	logrus.Error(err)
	// }

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	//ping database to make sure connection is established successfully
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	logrus.Info("Connect database successfully")
	return db, err
}

// Migrate - Database
func Migrate(db *gorm.DB) error {
	logrus.Info("running database migration")

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := mysqlMigration.WithInstance(sqlDB, &mysqlMigration.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"mysql", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err == migrate.ErrNoChange {
		logrus.Info("no schema changes to apply")
		return nil
	}

	return err
}
