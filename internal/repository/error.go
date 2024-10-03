package repository

import (
	"net/http"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
)

// Logs any Gorm database errors and will convert to 500 status code errors to be returned to the client.
func LogDBErr(logger *applog.AppLogger, gormErr error, msg string) error {
	if gormErr == nil {
		return nil
	}

	logger.Debug(msg, gormErr)
	return apperror.New(http.StatusInternalServerError, "Internal Server Error, check log for details")
}
