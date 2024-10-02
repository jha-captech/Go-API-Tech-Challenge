package apperror_test

import (
	"fmt"
	"net/http"
	"testing"

	"jf.go.techchallenge/internal/apperror"
)

func TestAppErrorNew(t *testing.T) {

	appErr := apperror.New(401, "Testing")

	if appErr.Message != "Testing" {
		t.Errorf("Want: %s got: %s", "Testing", appErr.Message)
	}

	if appErr.Status() != 401 {
		t.Errorf("Want: %d got: %d", 401, appErr.Status())
	}

	appErr = apperror.New(402, "Args %s %d", "foo", 5)

	if appErr.Message != "Args foo 5" {
		t.Errorf("Want: %s got: %s", "Args foo 5", appErr.Message)
	}

	if appErr.Status() != 402 {
		t.Errorf("Want: %d got: %d", 402, appErr.Status())
	}

}

func TestAppErrorBadRequest(t *testing.T) {
	appErr := apperror.BadRequest("Testing %s %d", "Foo", 10).(apperror.AppError)

	if appErr.Message != "Testing Foo 10" {
		t.Errorf("Want: %s got: %s", "Testing Foo 10", appErr.Message)
	}

	if appErr.Status() != 400 {
		t.Errorf("Want: %d got: %d", 400, appErr.Status())
	}
}

func TestAppErrorNotFound(t *testing.T) {
	appErr := apperror.NotFound("Testing %s %d", "Foo", 10).(apperror.AppError)

	if appErr.Message != "Testing Foo 10" {
		t.Errorf("Want: %s got: %s", "Testing Foo 10", appErr.Message)
	}

	if appErr.Status() != 404 {
		t.Errorf("Want: %d got: %d", 404, appErr.Status())
	}
}

func TestConvertStatusError(t *testing.T) {
	var appErr error
	appErr = apperror.NotFound("Testing %s %d", "Foo", 10).(apperror.AppError)

	converted := apperror.ConvertStatusError(appErr)

	if fmt.Sprint(converted) != fmt.Sprint(appErr) {
		t.Errorf("Error Conversion was not as expected")
	}

	err := fmt.Errorf("Error")
	appErr = apperror.ConvertStatusError(err)

	if fmt.Sprint(appErr) != fmt.Sprint(apperror.New(http.StatusInternalServerError, "Internal Server Error")) {
		t.Errorf("Non Http Status Error was not converted to 500 error.")
	}
}

var appErrorsOfTestCases = []struct {
	name           string
	inE            []error
	outErr         error
	expectedStatus int
}{
	{
		name:           "Len 0 returns nil",
		expectedStatus: 500,
	},

	{
		name:           "error returns internal server error",
		inE:            []error{fmt.Errorf("Error")},
		outErr:         apperror.New(http.StatusInternalServerError, "Internal Server Error"),
		expectedStatus: 500,
	},
	{
		name: "Returns worst status code",
		inE:  []error{apperror.New(200, "foo"), apperror.New(400, "bar")},
		outErr: apperror.MultiError{
			Message: "Multiple Errors:",
			Errors:  []apperror.StatusError{apperror.New(200, "foo"), apperror.New(400, "bar")},
		},
		expectedStatus: 400,
	},
}

func TestAppErrorsOf(t *testing.T) {
	for _, tc := range appErrorsOfTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := apperror.Of(tc.inE)
			if fmt.Sprint(err) != fmt.Sprint(tc.outErr) {
				t.Errorf("Of was not correct want: %v got: %v", tc.outErr, err)
			}

			if appErr := apperror.ConvertStatusError(err); tc.expectedStatus != appErr.Status() {
				t.Errorf("Incorrect Http Status want: %d got %d", tc.expectedStatus, appErr.Status())
			}
		})
	}
}
