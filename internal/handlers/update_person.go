package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type PersonUpdater interface {
	UpdatePerson(ctx context.Context, lastName string, updatedPerson models.Person) (models.Person, error)
}

// HandleUpdatePerson is a Handler that updates a given person
//
//	@Summary		Update Person
//	@Description	Updates person
//	@Tags			person
//	@Accept			json
//	@Produce		json
//	@Param			name				path		string	true "last name of person to update"
//	@Param			person				body		handlers.inputPerson	true	"Person Object"
//	@Success		200					{object}	handlers.responsePerson
//	@Failure		500					{object}	handlers.responseErr
//	@Router			/api/person/{name}	[PUT]
func HandleUpdatePerson(logger *httplog.Logger, service PersonUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()
		lastName := chi.URLParam(r, "name")

		// get values from database
		personIn, problems, err := decodeValidateBody[inputPerson, models.Person](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					ValidationErrors: problems,
				})
			default:
				logger.Error("BodyParser error", "error", err)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		person, err := service.UpdatePerson(ctx, lastName, personIn)
		if err != nil {
			logger.Error("error updating person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		personOut := mapOutputPerson(person)
		encodeResponse(w, logger, http.StatusOK, responsePerson{
			Person: personOut,
		})
	}
}
