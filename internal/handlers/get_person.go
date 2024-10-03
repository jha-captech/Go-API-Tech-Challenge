package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type PersonGetter interface {
	GetPersonByName(ctx context.Context, name string) (models.Person, error)
}

// GetPerson is a Handler that returns a specific person by name.
//
//	@Summary		Gets Person
//	@Description	Gets Person by Name
//	@Tags			person
//	@Accept			json
//	@Produce		json
//	@Param			name				path		string	true "last name of person to retrieve"
//	@Success		200					{object}	handlers.responsePerson
//	@Failure		500					{object}	handlers.responseErr
//	@Router			/api/person/{name}	[GET]
func HandleGetPersonByName(logger *httplog.Logger, service PersonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()
		name := chi.URLParam(r, "name")

		// get values from database
		persons, err := service.GetPersonByName(ctx, name)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		personOut := mapOutputPerson(persons)
		encodeResponse(w, logger, http.StatusOK, responsePerson{
			Person: personOut,
		})
	}
}
