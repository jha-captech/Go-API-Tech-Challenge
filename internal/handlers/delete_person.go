package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type PersonDeleter interface {
	DeletePerson(ctx context.Context, lastName string) error
}

// HandleDeletePerson is a Handler that deletes a given person
//
//	@Summary		Deletes Person
//	@Description	Deletes person by name
//	@Tags			person
//	@Accept			json
//	@Produce		json
//	@Param			name				path		string	true "last name of person to delete"
//	@Success		200					{object}	handlers.responseMsg
//	@Failure		500					{object}	handlers.responseErr
//	@Router			/api/person/{name}	[DELETE]
func HandleDeletePerson(logger *httplog.Logger, service PersonDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()
		lastName := chi.URLParam(r, "name")

		err := service.DeletePerson(ctx, lastName)
		if err != nil {
			logger.Error("error deleting person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error deleting person",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, responseMsg{
			Message: "Person deleted successfully",
		})
	}
}
