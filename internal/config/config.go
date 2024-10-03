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

type EnvProvider interface {
	Env(property string) (string, bool)
}

type OsEnvProvider struct{}

func (OsEnvProvider) Env(property string) (string, bool) {
	return os.LookupEnv(property)
}

// Creates and returns a new configuration based on environment variables passed to the executing program.
// Will return an error if any required environment variables are not set.
func New(log *applog.AppLogger) (*Configuration, error) {
	return NewWithProvider(log, OsEnvProvider{})
}

func NewWithProvider(log *applog.AppLogger, provider EnvProvider) (*Configuration, error) {
	databaseName, databaseNameSet := provider.Env("DATABASE_NAME")
	if !databaseNameSet {
		return nil, fmt.Errorf("database environment variable must be set")
	}

	databaseUser, databaseUserSet := provider.Env("DATABASE_USER")
	if !databaseUserSet {
		return nil, fmt.Errorf("database user environment variable must be set")
	}

	databasePassword, databasePasswordSet := provider.Env("DATABASE_PASSWORD")
	if !databasePasswordSet {
		return nil, fmt.Errorf("database passsord environment variable must be set")
	}

	databaseHost, databaseHostSet := provider.Env("DATABASE_HOST")
	if !databaseHostSet {
		return nil, fmt.Errorf("database host environment variable must bet set")
	}

	databasePort, databasePortSet := provider.Env("DATABASE_PORT")
	if !databasePortSet {
		return nil, fmt.Errorf("database port environment variable must be set")
	}

	retryString, retrySet := provider.Env("DATABASE_RETRY_DURATION_SECONDS")
	if !retrySet {
		retryString = "5"
	}

	databaseRetry, err := strconv.Atoi(retryString)
	if err != nil {
		log.Debug("Error converting DATABASE_RETRY_DURATION_SECONDS property", err)
		return nil, fmt.Errorf("DATABASE_RETRY_DURATION_SECONDS must be a number ")
	}

	logLevel, logLevelSet := provider.Env("LOG_LEVEL")
	if !logLevelSet {
		logLevel = "DEBUG"
	}

	config := &Configuration{
		LogLevel: logLevel,
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
