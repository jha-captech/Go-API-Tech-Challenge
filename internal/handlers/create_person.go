package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type PersonCreator interface {
	CreatePerson(ctx context.Context, newPerson models.Person) (models.Person, error)
}

// HandleCreatePerson is a Handler that creates a new person
//
//	@Summary		Creates Person
//	@Description	Creates person
//	@Tags			person
//	@Accept			json
//	@Produce		json
//	@Param			person		body		handlers.inputPerson	true	"Person Object"
//	@Success		200			{object}	handlers.responsePerson
//	@Failure		500			{object}	handlers.responseErr
//	@Router			/api/person	[POST]
func HandleCreatePerson(logger *httplog.Logger, service PersonCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from request body
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

		person, err := service.CreatePerson(ctx, personIn)
		if err != nil {
			logger.Error("error creating person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error creating person",
			})
			return
		}

		personOut := mapOutputPerson(person)
		encodeResponse(w, logger, http.StatusCreated, responsePerson{
			Person: personOut,
		})
	}
}
