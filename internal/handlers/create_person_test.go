package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	serviceMock "go-api-tech-challenge/internal/handlers/mock"
	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreatePerson(t *testing.T) {
	mockService := new(serviceMock.PersonCreator)
	logger := httplog.NewLogger("test")
	handler := HandleCreatePerson(logger, mockService)

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
		body         string
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"person created successfully": {
			body:         `{"first_name": "John", "last_name": "Doe", "type": "student", "age": 25, "courses": [1, 2]}`,
			mockCalled:   true,
			mockOutput:   []any{personOut, nil},
			expectedCode: http.StatusCreated,
			expectedBody: testutil.ToJSONString(responsePerson{Person: mapOutputPerson(personOut)}),
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
					{Name: "first_name", Description: "must not be blank"},
					{Name: "last_name", Description: "must not be blank"},
					{Name: "type", Description: "must be either 'student' or 'professor'"},
				},
			}),
		},
		"internal server error": {
			body:         `{"first_name": "John", "last_name": "Doe", "type": "student", "age": 25, "courses": [1, 2]}`,
			mockCalled:   true,
			mockOutput:   []any{models.Person{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error creating person"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/person", strings.NewReader(tc.body))
			assert.NoError(t, err)

			if tc.mockCalled {
				mockService.
					On("CreatePerson", mock.Anything, mock.MatchedBy(func(p models.Person) bool {
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
				mockService.AssertNotCalled(t, "CreatePerson")
			}
		})
	}
}
