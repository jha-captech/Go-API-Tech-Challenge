package handler

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
)

type MockHttpWriter struct {
	mock.Mock
}

func (m MockHttpWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m MockHttpWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

func (m MockHttpWriter) Write(bytes []byte) (int, error) {
	args := m.Called(bytes)
	return args.Int(0), args.Error(1)
}

func TestEncodeErrorExpectingStatusFromAppError(t *testing.T) {
	mockHttpWriter := new(MockHttpWriter)

	mockHttpWriter.On("WriteHeader", 404)

	mockHttpWriter.On("Write", mock.Anything).Return(1, nil)

	encodeError(mockHttpWriter, apperror.NotFound("Not Found"))

	mockHttpWriter.AssertExpectations(t)
}

func TestEncodeErrorExpectingStatusFromConvert(t *testing.T) {
	mockHttpWriter := new(MockHttpWriter)

	mockHttpWriter.On("WriteHeader", 500)

	mockHttpWriter.On("Write", mock.Anything).Return(1, nil)

	encodeError(mockHttpWriter, fmt.Errorf("NullPointerException"))

	mockHttpWriter.AssertExpectations(t)
}

var encodeResponseTestCases = []struct {
	name string

	onWriteInt   int
	onWriteError error

	passedError error

	expectedStatus int
}{
	{
		name:           "Success",
		onWriteInt:     1,
		onWriteError:   nil,
		expectedStatus: 200,
		passedError:    nil,
	},
	{
		name:           "With Error",
		onWriteInt:     1,
		onWriteError:   nil,
		expectedStatus: 400,
		passedError:    apperror.BadRequest("No Bueno"),
	},
	{
		name:           "Encoding Failure",
		onWriteInt:     -1,
		onWriteError:   fmt.Errorf("NullPointerException"),
		expectedStatus: 500,
		passedError:    nil,
	},
}

func TestEncodeResponse(t *testing.T) {

	applog := applog.New(log.Default())

	for _, tc := range encodeResponseTestCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHttpWriter := new(MockHttpWriter)

			mockHttpWriter.On("Header").Return(http.Header{})
			mockHttpWriter.On("Write", mock.Anything).Return(tc.onWriteInt, tc.onWriteError)
			if tc.expectedStatus != 200 {
				mockHttpWriter.On("WriteHeader", tc.expectedStatus)
			}
			encodeResponse(mockHttpWriter, applog, "Hello!", tc.passedError)

			mockHttpWriter.AssertExpectations(t)
		})
	}
}
