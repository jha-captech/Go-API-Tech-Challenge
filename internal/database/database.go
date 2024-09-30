package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jf.go.techchallenge/internal/config"
)

// Start Database
func New(config *config.Configuration, goLog *log.Logger) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)

	// Gorm Logger
	newLogger := logger.New(
		goLog, // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,
		},
	)

	// Connect to database.
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{
		// Logger: newLogger, todo
		Logger: newLogger,
	})
}
