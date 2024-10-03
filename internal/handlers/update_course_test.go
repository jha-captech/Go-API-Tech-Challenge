package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdateCourse(t *testing.T) {
	mockService := new(serviceMock.CourseUpdater)
	logger := httplog.NewLogger("test")
	handler := HandleUpdateCourse(logger, mockService)

	course := models.Course{ID: 1, Name: "Databases"}
	courseOut := mapOutputCourse(course)

	tests := map[string]struct {
		courseID     string
		body         string
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"course updated successfully": {
			courseID:     "1",
			body:         `{"name": "Databases"}`,
			mockCalled:   true,
			mockOutput:   []any{course, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseCourse{Course: courseOut}),
		},
		"invalid course ID": {
			courseID:     "abc",
			body:         `{"name": "Databases"}`,
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Not a valid ID"}),
		},
		"invalid body": {
			courseID:     "1",
			body:         `invalid body`,
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "missing values or malformed body"}),
		},
		"validation errors in body": {
			courseID:     "1",
			body:         `{}`,
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{
				ValidationErrors: []problem{
					{Name: "name", Description: "must not be blank"},
				},
			}),
		},
		"internal server error": {
			courseID:     "1",
			body:         `{"name": "Databases"}`,
			mockCalled:   true,
			mockOutput:   []any{models.Course{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, "/api/course/"+tc.courseID, strings.NewReader(tc.body))
			assert.NoError(t, err)

			// Add chi URLParam for courseID
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.courseID)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				id, _ := strconv.Atoi(tc.courseID) // Convert courseID to integer
				mockService.
					On("UpdateCourse", ctx, id, course.Name).
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
				mockService.AssertNotCalled(t, "UpdateCourse")
			}
		})
	}
}
