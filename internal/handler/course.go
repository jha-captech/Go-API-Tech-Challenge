package handler

import (
	"net/http"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/services"
)

func GetAllCourses(service *services.Course, logger *applog.AppLogger) Route {
	return NewRoute("GET /api/course", func(w http.ResponseWriter, r *http.Request) {
		courses, err := service.GetAll(r.URL.Query())
		encodeResponse(w, logger, courses, err)
	})
}

func CreateCourse(service *services.Course, logger *applog.AppLogger) Route {
	return NewRoute("POST /api/course", func(w http.ResponseWriter, r *http.Request) {
		input, err := decodeBody[models.CourseInput](r)

		if err != nil {
			logger.Debug("Failed to Deserialize Request Body ", err)
			encodeError(w, apperror.BadRequest("Invalid JSON"))
			return
		}
		course, err := service.Create(*input)
		encodeCreated(w, logger, course, err)
	})
}

func GetOneCourse(service *services.Course, logger *applog.AppLogger) Route {
	return NewRoute("GET /api/course/{guid}", func(w http.ResponseWriter, r *http.Request) {
		course, err := service.GetOneByGuid(r.PathValue("guid"))
		encodeResponse(w, logger, course, err)
	})
}

func UpdateOneCourse(service *services.Course, logger *applog.AppLogger) Route {
	return NewRoute("PUT /api/course/{guid}", func(w http.ResponseWriter, r *http.Request) {

		input, err := decodeBody[models.CourseInput](r)

		if err != nil {
			logger.Debug("Failed to Deserialize Request Body ", err)
			encodeError(w, apperror.BadRequest("Invalid JSON"))
			return
		}

		course, err := service.Update(r.PathValue("guid"), *input)
		encodeResponse(w, logger, course, err)
	})
}

func DeleteOneCourse(service *services.Course, logger *applog.AppLogger) Route {
	return NewRoute("DELETE /api/course/{guid}", func(w http.ResponseWriter, r *http.Request) {
		err := service.Delete(r.PathValue("guid"))
		encodeResponse(w, logger, "OK", err)
	})
}
