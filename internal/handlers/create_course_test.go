package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"

	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateCourse(t *testing.T) {
	mockService := new(serviceMock.CourseCreator)
	logger := httplog.NewLogger("test")
	handler := HandleCreateCourse(logger, mockService)

	course := models.Course{ID: 1, Name: "Databases"}
	courseOut := mapOutputCourse(course)

	tests := map[string]struct {
		body         string
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"course created successfully": {
			body:         `{"name": "Databases"}`,
			mockCalled:   true,
			mockOutput:   []any{course, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseCourse{Course: courseOut}),
		},
		"invalid body": {
			body:         `invalid body`,
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "missing values or malformed body"}),
		},
		"validation errors in body": {
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
			body:         `{"name": "Databases"}`,
			mockCalled:   true,
			mockOutput:   []any{models.Course{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/course", strings.NewReader(tc.body))
			assert.NoError(t, err)

			if tc.mockCalled {
				mockService.
					On("CreateCourse", context.Background(), "Databases").
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
				mockService.AssertNotCalled(t, "CreateCourse")
			}
		})
	}
}
