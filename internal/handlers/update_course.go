package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/httplog/v2"
)

type CourseUpdater interface {
	UpdateCourse(ctx context.Context, courseID int, courseName string) (models.Course, error)
}

// UpdateCourse is a Handler that returns the course associated with the given ID.
//
//	@Summary		Update Course
//	@Description	Updates course associated with given ID
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			ID					path		int	true "ID of course to update"
//	@Param			course				body		handlers.inputCourse	true	"Course Object"
//	@Success		200					{object}	handlers.responseCourse
//	@Failure		500					{object}	handlers.responseErr
//	@Router			/api/course/{ID}	[PUT]
func HandleUpdateCourse(logger *httplog.Logger, service CourseUpdater) http.HandlerFunc {
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

		// get and validate body as object
		courseIn, problems, err := decodeValidateBody[inputCourse, models.Course](r)
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
		// get values from database
		course, err := service.UpdateCourse(ctx, courseID, courseIn.Name)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		courseOut := mapOutputCourse(course)
		encodeResponse(w, logger, http.StatusOK, responseCourse{
			Course: courseOut,
		})
	}
}
