package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type PersonLister interface {
	ListPersons(ctx context.Context) ([]models.Person, error)
}

// HandleListPersons is a Handler that returns a list of all persons.
//
//	@Summary		List all Persons
//	@Description	List all persons
//	@Tags			person
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	handlers.responsePersons
//	@Failure		500			{object}	handlers.responseErr
//	@Router			/api/person	[GET]
func HandleListPersons(logger *httplog.Logger, service PersonLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		persons, err := service.ListPersons(ctx)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		personsOut := mapMultipleOutputPerson(persons)
		encodeResponse(w, logger, http.StatusOK, responsePersons{
			Persons: personsOut,
		})
	}
}
