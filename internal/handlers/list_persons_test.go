package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"
	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleListPersons(t *testing.T) {
	mockService := new(serviceMock.PersonLister)
	logger := httplog.NewLogger("test")
	handler := HandleListPersons(logger, mockService)

	persons := []models.Person{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Type:      "student",
			Age:       25,
			Courses:   []int{1, 2},
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Type:      "professor",
			Age:       35,
			Courses:   []int{1},
		},
	}

	personsOut := mapMultipleOutputPerson(persons)

	tests := map[string]struct {
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"persons found": {
			mockCalled:   true,
			mockOutput:   []any{persons, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responsePersons{Persons: personsOut}),
		},
		"no persons found": {
			mockCalled:   true,
			mockOutput:   []any{[]models.Person{}, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responsePersons{Persons: []outputPerson{}}),
		},
		"internal server error": {
			mockCalled:   true,
			mockOutput:   []any{[]models.Person{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/persons", nil)
			assert.NoError(t, err)

			if tc.mockCalled {
				mockService.
					On("ListPersons", context.Background()).
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
				mockService.AssertNotCalled(t, "ListPersons")
			}
		})
	}
}
