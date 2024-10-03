package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/httplog/v2"
)

type CourseDeleter interface {
	DeleteCourse(ctx context.Context, courseID int) error
}

// DeleteCourse is a Handler that returns the course associated with the given ID.
//
//	@Summary		Deletes Course
//	@Description	Deletes course associated with given ID
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			ID					path	int		true "ID of course to delete"
//	@Success		200					{object}	handlers.responseMsg
//	@Failure		500					{object}	handlers.responseErr
//	@Router			/api/course/{ID}	[DELETE]
func HandleDeleteCourse(logger *httplog.Logger, service CourseDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// setup
		idString := chi.URLParam(r, "ID")
		courseID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
			return
		}

		err = service.DeleteCourse(ctx, courseID)
		if err != nil {

			logger.Error("error deleting course", "error", err)

			if err.Error() == "course not found" {
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "course ID not found",
				})
				return
			}
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, responseMsg{
			Message: "Course deleted successfully",
		})
	}
}
