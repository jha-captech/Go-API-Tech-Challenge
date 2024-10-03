package handlers

import (
	"context"
	"go-api-tech-challenge/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/httplog/v2"
)

type CourseGetter interface {
	GetCourseByID(ctx context.Context, ID int) (models.Course, error)
}

// GetCourseByID is a Handler that returns the course associated with the given ID.
//
//	@Summary		Get Course
//	@Description	Gets course associated with given ID
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			ID					path		int	true "ID of course to retrieve"
//	@Success		200					{object}	handlers.responseCourse
//	@Failure		500					{object}	handlers.responseErr
//	@Router			/api/course/{ID}	[GET]
func HandleGetCourseByID(logger *httplog.Logger, service CourseGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// setup
		idString := chi.URLParam(r, "ID")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Error retrieving course",
			})
			return
		}
		// get values from database
		course, err := service.GetCourseByID(ctx, ID)
		if err != nil {
			log.Println(err)
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
