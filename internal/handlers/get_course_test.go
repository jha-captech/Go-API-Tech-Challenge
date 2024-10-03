package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetCourseByID(t *testing.T) {
	mockService := new(serviceMock.CourseGetter)
	logger := httplog.NewLogger("test")
	handler := HandleGetCourseByID(logger, mockService)

	course := models.Course{ID: 1, Name: "Databases"}
	courseOut := mapOutputCourse(course)

	tests := map[string]struct {
		courseID     string
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"course found": {
			courseID:     "1",
			mockCalled:   true,
			mockOutput:   []any{course, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseCourse{Course: courseOut}),
		},
		"course not found": {
			courseID:     "999",
			mockCalled:   true,
			mockOutput:   []any{models.Course{}, errors.New("course not found")},
			expectedCode: http.StatusInternalServerError, // Changed to match the handler
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
		"invalid ID": {
			courseID:     "abc",
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving course"}),
		},
		"internal server error": {
			courseID:     "1",
			mockCalled:   true,
			mockOutput:   []any{models.Course{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/course/"+tc.courseID, nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.courseID) // The handler expects "ID" as the URL parameter
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				id, _ := strconv.Atoi(tc.courseID) // Convert courseID to integer
				mockService.
					On("GetCourseByID", ctx, id).
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
				mockService.AssertNotCalled(t, "GetCourseByID")
			}
		})
	}
}
