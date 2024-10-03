package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleListUsers(t *testing.T) {
	mockService := new(serviceMock.CourseLister)
	logger := httplog.NewLogger("test")
	handler := HandleListCourses(logger, mockService)

	courses := []models.Course{
		{ID: 1, Name: "Databases"},
		{ID: 2, Name: "Operating Systems"},
	}

	coursesOut := mapMultipleOutputCourse(courses)

	tests := map[string]struct {
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"courses returned": {
			mockCalled:   true,
			mockOutput:   []any{courses, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseCourses{Courses: coursesOut}),
		},
		"no users found": {
			mockCalled:   true,
			mockOutput:   []any{[]models.Course{}, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseCourses{Courses: []outputCourse{}}),
		},
		"internal server error": {
			mockCalled:   true,
			mockOutput:   []any{[]models.Course{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/course", nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("ListCourses", ctx).
					Return(tc.mockOutput...).
					Once()
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "ListCourses")
			}
		})
	}
}
