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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetPersonByName(t *testing.T) {
	mockService := new(serviceMock.PersonGetter)
	logger := httplog.NewLogger("test")
	handler := HandleGetPersonByName(logger, mockService)

	person := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	personOut := mapOutputPerson(person)

	tests := map[string]struct {
		name         string
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"person found": {
			name:         "Doe",
			mockCalled:   true,
			mockOutput:   []any{person, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responsePerson{Person: personOut}),
		},
		"person not found": {
			name:         "Smith",
			mockCalled:   true,
			mockOutput:   []any{models.Person{}, errors.New("person not found")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
		"internal server error": {
			name:         "Doe",
			mockCalled:   true,
			mockOutput:   []any{models.Person{}, errors.New("test error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/person/"+tc.name, nil)
			assert.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", tc.name)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("GetPersonByName", ctx, tc.name).
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
				mockService.AssertNotCalled(t, "GetPersonByName")
			}
		})
	}
}
