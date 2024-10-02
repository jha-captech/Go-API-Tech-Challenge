package handler_test

import (
	"net/http"

	"github.com/stretchr/testify/mock"
	"github.com/vektra/mockery/mockery/fixtures/mocks"
)

func TestEncodeErrorExpectingStatusErr() {
	mockResponseWriter := new(mocks.ResponseWriter)

	mockResponseWriter.On("WriteHeader", http.NotFound).Once()
	mockResponseWriter.On("Write", mock.Anything).Return(len([]byte("Hello, World!")), nil)

}
