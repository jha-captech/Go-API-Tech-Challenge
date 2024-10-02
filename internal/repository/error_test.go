package repository_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/repository"
)

func TestLogDBErr(t *testing.T) {
	logger := applog.New(log.Default())

	err := repository.LogDBErr(logger, fmt.Errorf("Err"), "This is an error!")

	if fmt.Sprint(err) != fmt.Sprint(apperror.New(http.StatusInternalServerError, "Internal Server Error, check log for details")) {
		t.Errorf("Error was not correct")
	}
}
