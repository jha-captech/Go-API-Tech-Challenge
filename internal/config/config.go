package config

import (
	"fmt"
	"os"
	"strconv"

	"jf.go.techchallenge/internal/applog"
)

type Database struct {
	User         string
	Password     string
	Name         string
	Host         string
	Port         string
	RetrySeconds int
}

type Configuration struct {
	Database Database
	LogLevel string
}

// Creates and returns a new configuration based on environment variables passed to the executing program.
// Will return an error if any required environment variables are not set.
func New(log *applog.AppLogger) (*Configuration, error) {
	databaseName, databaseNameSet := os.LookupEnv("DATABASE_NAME")
	if !databaseNameSet {
		return nil, fmt.Errorf("database environment variable must be set")
	}

	databaseUser, databaseUserSet := os.LookupEnv("DATABASE_USER")
	if !databaseUserSet {
		return nil, fmt.Errorf("database user environment variable must be set")
	}

	databasePassword, databasePasswordSet := os.LookupEnv("DATABASE_PASSWORD")
	if !databasePasswordSet {
		return nil, fmt.Errorf("database passsord environment variable must be set")
	}

	databaseHost, databaseHostSet := os.LookupEnv("DATABASE_HOST")
	if !databaseHostSet {
		return nil, fmt.Errorf("database host environment variable must bet set")
	}

	databasePort, databasePortSet := os.LookupEnv("DATABASE_PORT")
	if !databasePortSet {
		return nil, fmt.Errorf("database port environment variable must be set")
	}

	retryString, retrySet := os.LookupEnv("DATABASE_RETRY_DURATION_SECONDS")
	if !retrySet {
		retryString = "5"
	}
	databaseRetry, err := strconv.Atoi(retryString)
	if err != nil {
		log.Fatal("DATABASE_RETRY_DURATION_SECONDS must be a number ", err)
	}

	config := &Configuration{
		Database: Database{
			RetrySeconds: databaseRetry,
			User:         databaseUser,
			Name:         databaseName,
			Password:     databasePassword,
			Host:         databaseHost,
			Port:         databasePort,
		},
	}

	return config, nil
}
