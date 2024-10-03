package config

import (
	"fmt"
	"os"
)

type Database struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

type Configuration struct {
	Database Database
}

// Creates and returns a new configuration based on environment variables passed to the executing program.
// Will return an error if any required environment variables are not set.
func New() (*Configuration, error) {
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

	config := &Configuration{
		Database: Database{
			User:     databaseUser,
			Name:     databaseName,
			Password: databasePassword,
			Host:     databaseHost,
			Port:     databasePort,
		},
	}

	return config, nil
}
