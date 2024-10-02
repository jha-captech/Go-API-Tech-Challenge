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

// type MockHttpResponseWriter struct {
// 	expectStatus     int
// 	wasWriteCalled   bool
// 	statusCalledWith int
// 	shouldWriteError error
// }

// func (s MockHttpResponseWriter) Header() http.Header {
// 	return http.Header{}
// }

// func (s MockHttpResponseWriter) Write([]byte) (int, error) {
// 	s.wasWriteCalled = true
// 	if s.shouldWriteError != nil {
// 		return -1, s.shouldWriteError
// 	}
// 	return 1, nil
// }

// func (s MockHttpResponseWriter) WriteHeader(statusCode int) {
// 	s.statusCalledWith = statusCode
// }

// func (s MockHttpResponseWriter) GetStatusCalled() int {
// 	return s.statusCalledWith
// }

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

// func TestEncodeResponseExpectingResponseEncoded(t *testing.T) {
// 	mockHttpWriter := new(MockHttpWriter)

// 	applog := applog.AppLogger{}

// 	mockHttpWriter.On("Header").Return(http.Header{})

// 	mockHttpWriter.On("Write", mock.Anything).Return(1, nil)

// 	encodeResponse(mockHttpWriter, &applog, "Hello!", nil)

// 	mockHttpWriter.AssertExpectations(t)
// }

// func TestEncodeResponseWithError(t *testing.T) {
// 	mockHttpWriter := new(MockHttpWriter)

// 	applog := applog.New(log.Default())

// 	mockHttpWriter.On("Header").Return(http.Header{})

// 	mockHttpWriter.On("Write", mock.Anything).Return(1, nil)

// 	mockHttpWriter.On("WriteHeader", 400)

// 	encodeResponse(mockHttpWriter, applog, "Hello!", apperror.BadRequest("No Bueno"))

// 	mockHttpWriter.AssertExpectations(t)
// }
