package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"
	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUpdatePerson(t *testing.T) {
	mockService := new(serviceMock.PersonUpdater)
	logger := httplog.NewLogger("test")
	handler := HandleUpdatePerson(logger, mockService)

	personIn := models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	personOut := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	tests := map[string]struct {
		lastName     string
		body         string
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"person updated successfully": {
			lastName:     "Doe",
			body:         `{"first_name": "John", "last_name": "Doe", "type": "student", "age": 25, "courses": [1, 2]}`,
			mockCalled:   true,
			mockOutput:   []any{personOut, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responsePerson{Person: mapOutputPerson(personOut)}),
		},
		"invalid body": {
			lastName:     "Doe",
			body:         `invalid body`,
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{Error: "missing values or malformed body"}),
		},
		"validation errors in body": {
			lastName:     "Doe",
			body:         `{}`,
			mockCalled:   false,
			expectedCode: http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{
				ValidationErrors: []problem{
					{Name: "first_name", Description: "must not be blank"},
					{Name: "last_name", Description: "must not be blank"},
					{Name: "type", Description: "must be either 'student' or 'professor'"},
				},
			}),
		},
		"person not found": {
			lastName:     "Smith",
			body:         `{"first_name": "John", "last_name": "Doe", "type": "student", "age": 25, "courses": [1, 2]}`,
			mockCalled:   true,
			mockOutput:   []any{models.Person{}, errors.New("person not found")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
		"internal server error": {
			lastName:     "Doe",
			body:         `{"first_name": "John", "last_name": "Doe", "type": "student", "age": 25, "courses": [1, 2]}`,
			mockCalled:   true,
			mockOutput:   []any{models.Person{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, "/api/person/"+tc.lastName, strings.NewReader(tc.body))
			assert.NoError(t, err)

			// Add chi URLParam for last name
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", tc.lastName)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("UpdatePerson", mock.Anything, tc.lastName, mock.MatchedBy(func(p models.Person) bool {
						return p.FirstName == personIn.FirstName &&
							p.LastName == personIn.LastName &&
							p.Type == personIn.Type &&
							p.Age == personIn.Age &&
							len(p.Courses) == len(personIn.Courses) &&
							p.Courses[0] == personIn.Courses[0] &&
							p.Courses[1] == personIn.Courses[1]
					})).
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
				mockService.AssertNotCalled(t, "UpdatePerson")
			}
		})
	}
}
