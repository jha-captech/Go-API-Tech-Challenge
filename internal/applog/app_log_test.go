package applog_test

import (
	"fmt"
	"log"
	"testing"

	"jf.go.techchallenge/internal/applog"
)

func TestLogging(t *testing.T) {
	appLog := applog.New(log.Default())

	appLog.Debug("This is a debug message %d %s", 5, "hello!")
	appLog.Info("This is an info message %d %s", 5, "hello!")
	appLog.Error("This is an error message %d %s", fmt.Errorf("Error!"))
}
