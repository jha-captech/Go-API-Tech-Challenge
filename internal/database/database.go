package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/config"
)

// Start Database
func New(config *config.Configuration, appLogger *applog.AppLogger) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)

	level := logger.Silent
	// var
	if config.LogLevel == "DEBUG" {
		level = logger.Info
	}

	// Gorm Logger
	newLogger := logger.New(
		appLogger.GoLogger(), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  level,       // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,
		},
	)

	for attempts := 0; attempts < 30; attempts++ {
		db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
			Logger: newLogger,
		})

		if err == nil {
			return db, nil
		}

		appLogger.Info("Attempt %d: Unable to connect to the database: %v", attempts+1, err)
		time.Sleep(time.Duration(config.Database.RetrySeconds) * time.Second)
	}

	return nil, fmt.Errorf("FATAL: Failed to Connect to Database Tried 30 times and gave up")
}
