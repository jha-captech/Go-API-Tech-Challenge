package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type CourseCreator interface {
	CreateCourse(ctx context.Context, courseName string) (models.Course, error)
}

// CreateCourse is a Handler that creates a new course
//
//	@Summary		Create Course
//	@Description	Creates new course
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			course		body		handlers.inputCourse	true	"Course Object"
//	@Success		200			{object}	handlers.responseCourse
//	@Failure		500			{object}	handlers.responseErr
//	@Router			/api/course	[POST]
func HandleCreateCourse(logger *httplog.Logger, service CourseCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// get and validate body as object
		courseIn, problems, err := decodeValidateBody[inputCourse](r)
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
		course, err := service.CreateCourse(ctx, courseIn.Name)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		coursesOut := mapOutputCourse(course)
		encodeResponse(w, logger, http.StatusOK, responseCourse{
			Course: coursesOut,
		})
	}
}
