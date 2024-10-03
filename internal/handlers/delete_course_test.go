package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"
	"go-api-tech-challenge/internal/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleDeleteCourse(t *testing.T) {
	mockService := new(serviceMock.CourseDeleter)
	logger := httplog.NewLogger("test")
	handler := HandleDeleteCourse(logger, mockService)

	tests := map[string]struct {
		courseID     string
		mockCalled   bool
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		"course deleted successfully": {
			courseID:     "1",
			mockCalled:   true,
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseMsg{
				Message: "Course deleted successfully",
			}),
		},
		"invalid course ID": {
			courseID:     "abc",
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Not a valid ID"}),
		},
		"course not found": {
			courseID:     "1",
			mockCalled:   true,
			mockReturn:   errors.New("course not found"),
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "course ID not found"}),
		},
		"internal server error": {
			courseID:     "1",
			mockCalled:   true,
			mockReturn:   errors.New("test error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "/api/course/"+tc.courseID, nil)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.courseID)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				id, _ := strconv.Atoi(tc.courseID) // Convert courseID to integer
				mockService.
					On("DeleteCourse", ctx, id).
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
				mockService.AssertNotCalled(t, "DeleteCourse")
			}
		})
	}
}
