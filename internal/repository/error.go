package repository

import (
	"net/http"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
)

func logDBErr(logger *applog.AppLogger, gormErr error, msg string) error {
	if gormErr == nil {
		return nil
	}

	logger.Debug(msg, gormErr)
	return apperror.New(http.StatusInternalServerError, "Internal Server Error, check log for details")
}
