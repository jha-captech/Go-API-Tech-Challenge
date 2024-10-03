package config_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/mock"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/config"
)

type MockEnvProvider struct {
	mock.Mock
}

func (m MockEnvProvider) Env(property string) (string, bool) {
	args := m.Called(property)
	return args.Get(0).(string), args.Get(1).(bool)
}

var configTestCases = []struct {
	name        string
	configs     map[string]string
	expectedErr error
}{
	{
		name: "Success",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "10",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: nil,
	},
	{
		name: "Success DATABASE_RETRY_DURATION_SECONDS not required",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: nil,
	},
	{
		name: "Success LOG_LEVEL not required",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "5",
			"LOG_LEVEL":                       "",
		},
		expectedErr: nil,
	},
	{
		name: "DATABASE_NAME not set",
		configs: map[string]string{
			"DATABASE_NAME":                   "",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "10",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: fmt.Errorf("database environment variable must be set"),
	},
	{
		name: "DATABASE_USER not set",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "10",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: fmt.Errorf("database user environment variable must be set"),
	},
	{
		name: "DATABASE_HOST not set",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "10",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: fmt.Errorf("database host environment variable must bet set"),
	},
	{
		name: "DATABASE_PORT not set",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "",
			"DATABASE_RETRY_DURATION_SECONDS": "10",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: fmt.Errorf("database port environment variable must be set"),
	},
	{
		name: "DATABASE_RETRY_DURATION_SECONDS not a number",
		configs: map[string]string{
			"DATABASE_NAME":                   "test",
			"DATABASE_USER":                   "root",
			"DATABASE_PASSWORD":               "password",
			"DATABASE_HOST":                   "localhost",
			"DATABASE_PORT":                   "5432",
			"DATABASE_RETRY_DURATION_SECONDS": "a",
			"LOG_LEVEL":                       "DEBUG",
		},
		expectedErr: fmt.Errorf("DATABASE_RETRY_DURATION_SECONDS must be a number "),
	},
}

func TestConfig(t *testing.T) {
	logger := applog.New(log.Default())
	for _, tc := range configTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks.
			mockEnvProvider := new(MockEnvProvider)

			for key, value := range tc.configs {
				mockEnvProvider.On("Env", key).Return(value, value != "")
			}

			outConfig, outErr := config.NewWithProvider(logger, mockEnvProvider)

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Error was not as expected want: %v got: %v", tc.expectedErr, outErr)
			}

			if outErr != nil {
				return
			}

			if value := tc.configs["DATABASE_NAME"]; value != outConfig.Database.Name {
				t.Errorf("DATABASE_NAME was not correct %v", outConfig)
			}

			if value := tc.configs["DATABASE_USER"]; value != outConfig.Database.User {
				t.Errorf("DATABASE_USER was not correct %v", outConfig)
			}

			if value := tc.configs["DATABASE_PASSWORD"]; value != outConfig.Database.Password {
				t.Errorf("DATABASE_PASSWORD was not correct %v", outConfig)
			}

			if value := tc.configs["DATABASE_HOST"]; value != outConfig.Database.Host {
				t.Errorf("DATABASE_HOST was not correct %v", outConfig)
			}

			if value := tc.configs["DATABASE_PORT"]; value != outConfig.Database.Port {
				t.Errorf("DATABASE_PORT was not correct %v", outConfig)
			}

			if value := tc.configs["DATABASE_RETRY_DURATION_SECONDS"]; value != "" && value != fmt.Sprint(outConfig.Database.RetrySeconds) {
				t.Errorf("DATABASE_RETRY_DURATION_SECONDS was not correct %v", outConfig)
			}
			if value := tc.configs["LOG_LEVEL"]; value != "" && value != outConfig.LogLevel {
				t.Errorf("LOG_LEVEL was not correct %v", outConfig)
			}

		})
	}
}
