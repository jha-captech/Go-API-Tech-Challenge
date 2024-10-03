package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"
	"go-api-tech-challenge/internal/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleDeletePerson(t *testing.T) {
	mockService := new(serviceMock.PersonDeleter)
	logger := httplog.NewLogger("test")
	handler := HandleDeletePerson(logger, mockService)

	tests := map[string]struct {
		lastName     string
		mockCalled   bool
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		"person deleted successfully": {
			lastName:     "Doe",
			mockCalled:   true,
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseMsg{Message: "Person deleted successfully"}),
		},
		"person not found": {
			lastName:     "Smith",
			mockCalled:   true,
			mockReturn:   errors.New("person not found"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error deleting person"}),
		},
		"internal server error": {
			lastName:     "Doe",
			mockCalled:   true,
			mockReturn:   errors.New("test error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error deleting person"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "/api/person/"+tc.lastName, nil)
			assert.NoError(t, err)

			// Add chi URLParam for last name
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", tc.lastName)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("DeletePerson", mock.Anything, tc.lastName).
					Return(tc.mockReturn).
					Once()
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "DeletePerson")
			}
		})
	}
}
